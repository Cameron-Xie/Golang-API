package postgres

import (
	"fmt"

	"github.com/Cameron-Xie/Golang-API/pkg/http/rest"
	"github.com/Cameron-Xie/Golang-API/pkg/services/readtask"
	"github.com/Cameron-Xie/Golang-API/pkg/services/storetask"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

const taskTableName = "tasks"

type taskRow struct {
	readtask.Task
	Total int
}

func (s *Storage) FetchTasks(offset, limit int) (*rest.Collection, error) {
	sql := fmt.Sprintf(`
SELECT *
FROM (SELECT * FROM %v OFFSET ? LIMIT ?) i
         RIGHT JOIN (SELECT count(*) as total FROM %v) as c ON TRUE
`, taskTableName, taskTableName)

	tasks := make([]taskRow, 0)
	s.db.Table(taskTableName).Raw(sql, offset, limit).Scan(&tasks)

	return toTaskCollection(tasks), nil
}

func (s *Storage) FetchTask(id uuid.UUID) (*readtask.Task, error) {
	tasks := make([]readtask.Task, 0)
	s.db.Where("ID = ?", id.String()).Find(&tasks)

	if len(tasks) != 1 {
		return nil, &NotFoundError{
			Table: taskTableName,
			Value: id.String(),
		}
	}

	return &tasks[0], nil
}

func (s *Storage) StoreTask(i *storetask.Task) error {
	return s.db.Create(i).Error
}

// nolint interfacer
func (s *Storage) UpdateTask(i map[string]interface{}, id uuid.UUID) error {
	table := s.db.Table(taskTableName)
	if table.First(new(gorm.RowQueryResult), "id = ?", id.String()).RecordNotFound() {
		return &NotFoundError{
			Table: taskTableName,
			Value: id.String(),
		}
	}

	return table.Where("id = ?", id.String()).Update(i).Error
}

func (s *Storage) DeleteTask(id uuid.UUID) error {
	task := &storetask.Task{ID: id}
	if s.db.First(task).RecordNotFound() {
		return &NotFoundError{
			Table: taskTableName,
			Value: id.String(),
		}
	}

	return s.db.Delete(task).Error
}

func toTaskCollection(i []taskRow) *rest.Collection {
	if len(i) < 1 {
		panic("should return empty task with total count")
	}

	tasks := make([]readtask.Task, 0)
	for k := range i {
		if i[k].ID.String() != new(uuid.UUID).String() {
			tasks = append(tasks, i[k].Task)
		}
	}

	return &rest.Collection{
		Total: i[0].Total,
		Items: tasks,
	}
}
