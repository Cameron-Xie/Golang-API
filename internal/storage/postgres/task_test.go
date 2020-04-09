package postgres

import (
	"fmt"
	"testing"

	"github.com/Cameron-Xie/Golang-API/pkg/services/readtask"
	"github.com/Cameron-Xie/Golang-API/pkg/services/storetask"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestStorage_FetchTasks(t *testing.T) {
	a := assert.New(t)
	storage, err := New(getValidTestConn()).Open()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cleanUpTestData(storage.db)
		_ = storage.Close()
	}()

	// cleanup test data
	_ = cleanUpTestData(storage.db)

	// test with empty table
	if coll, err := storage.FetchTasks(0, 10); err == nil {
		tasks := coll.Items.([]readtask.Task)
		a.Equal(0, coll.Total)
		a.Equal(0, len(tasks))
	} else {
		t.Fatal(err)
	}

	// add task
	id := uuid.New()
	if err := storage.StoreTask(&storetask.Task{
		ID:          id,
		Name:        "task_a",
		Description: "description",
	}); err != nil {
		t.Fatal(err)
	}

	// test within range
	if coll, err := storage.FetchTasks(0, 10); err == nil {
		tasks := coll.Items.([]readtask.Task)
		a.Equal(1, coll.Total)
		a.Equal(1, len(tasks))
		a.Nil(
			isEqual(
				readtask.Task{
					ID:          id,
					Name:        "task_a",
					Description: "description",
				},
				tasks[0],
				[]string{"CreatedAt"},
			),
		)
	} else {
		t.Fatal(err)
	}

	// test out of range
	if coll, err := storage.FetchTasks(10, 5); err == nil {
		tasks := coll.Items.([]readtask.Task)
		a.Equal(1, coll.Total)
		a.Equal(0, len(tasks))
	} else {
		t.Fatal(err)
	}
}

func TestStorage_FetchTask(t *testing.T) {
	a := assert.New(t)
	storage, err := New(getValidTestConn()).Open()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cleanUpTestData(storage.db)
		_ = storage.Close()
	}()
	_ = cleanUpTestData(storage.db)

	// add task
	id, randomID := uuid.New(), uuid.New()
	if err := storage.StoreTask(&storetask.Task{
		ID:          id,
		Name:        "task_a",
		Description: "description",
	}); err != nil {
		t.Fatal(err)
	}

	m := []struct {
		id       uuid.UUID
		expected *readtask.Task
		err      error
	}{
		{
			id: id,
			expected: &readtask.Task{
				ID:          id,
				Name:        "task_a",
				Description: "description",
			},
		},
		{
			id: randomID,
			err: &NotFoundError{
				Table: taskTableName,
				Value: randomID.String(),
			},
		},
	}

	for _, i := range m {
		res, err := storage.FetchTask(i.id)
		a.Equal(i.err, err)
		if err == nil {
			a.Nil(isEqual(i.expected, res, []string{"CreatedAt"}))
			a.NotNil(res.CreatedAt)
			a.Nil(res.UpdatedAt)
		}
	}
}

func TestStorage_UpdateTask(t *testing.T) {
	a := assert.New(t)
	storage, err := New(getValidTestConn()).Open()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cleanUpTestData(storage.db)
		_ = storage.Close()
	}()
	_ = cleanUpTestData(storage.db)

	// add task
	id, randomID := uuid.New(), uuid.New()
	if err := storage.StoreTask(&storetask.Task{
		ID:          id,
		Name:        "task_a",
		Description: "description",
	}); err != nil {
		t.Fatal(err)
	}

	m := []struct {
		input    map[string]interface{}
		id       uuid.UUID
		expected *readtask.Task
		err      error
	}{
		{
			input: map[string]interface{}{
				"name": "new_task",
			},
			id: randomID,
			err: &NotFoundError{
				Table: taskTableName,
				Value: randomID.String(),
			},
		},
		{
			input: map[string]interface{}{
				"name": "new_task",
			},
			id: id,
			expected: &readtask.Task{
				ID:          id,
				Name:        "new_task",
				Description: "description",
			},
		},
	}

	for _, i := range m {
		err := storage.UpdateTask(i.input, i.id)

		if i.err != nil {
			a.Equal(i.err, err)
			continue
		}

		a.Nil(err)
		if res, err := storage.FetchTask(i.id); err == nil {
			a.Nil(isEqual(res, i.expected, []string{"CreatedAt", "UpdatedAt"}))
			a.NotNil(res.CreatedAt)
			a.NotNil(res.UpdatedAt)
		} else {
			t.Fatal(err)
		}
	}
}

func TestStorage_DeleteTask(t *testing.T) {
	a := assert.New(t)
	storage, err := New(getValidTestConn()).Open()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = cleanUpTestData(storage.db)
		_ = storage.Close()
	}()
	_ = cleanUpTestData(storage.db)

	// add task
	id, randomID := uuid.New(), uuid.New()
	if err := storage.StoreTask(&storetask.Task{
		ID:          id,
		Name:        "task_a",
		Description: "description",
	}); err != nil {
		t.Fatal(err)
	}

	m := []struct {
		id  uuid.UUID
		err error
	}{
		{
			id: id,
		},
		{
			id: randomID,
			err: &NotFoundError{
				Table: taskTableName,
				Value: randomID.String(),
			},
		},
	}

	for _, i := range m {
		a.Equal(i.err, storage.DeleteTask(i.id))
	}
}

func TestToTaskCollection(t *testing.T) {
	a := assert.New(t)
	a.Panics(func() {
		toTaskCollection(make([]taskRow, 0))
	})
}

func cleanUpTestData(db *gorm.DB) error {
	return db.Exec(
		fmt.Sprintf("DELETE FROM %v", taskTableName),
	).Error
}
