package validator

import (
	"errors"
	"fmt"
	"testing"

	pv "github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/stretchr/testify/assert"
)

func TestValidator_ValidateStruct(t *testing.T) {
	type testStruct struct {
		ID    string `validate:"required"`
		Email string `validate:"required,email"`
	}

	a := assert.New(t)
	v := &validator{v: pv.New()}
	m := []struct {
		s   *testStruct
		err error
	}{
		{
			s: &testStruct{},
			err: &InvalidParamsError{Errs: map[string]string{
				"email": "key: 'teststruct.email' error:field validation for 'email' failed on the 'required' tag",
				"id":    "key: 'teststruct.id' error:field validation for 'id' failed on the 'required' tag",
			}},
		},
		{
			s: &testStruct{Email: "example@test.com"},
			err: &InvalidParamsError{Errs: map[string]string{
				"id": "key: 'teststruct.id' error:field validation for 'id' failed on the 'required' tag",
			}},
		},
		{
			s: &testStruct{ID: "id", Email: "example@test.com"},
		},
	}

	for _, i := range m {
		a.Equal(i.err, v.ValidateStruct(i.s))
	}
}

func TestValidator_ValidateMap(t *testing.T) {
	a := assert.New(t)
	v := &validator{v: pv.New()}
	rules := map[string]MapRule{
		"id": {
			Rule: "required",
			Msg:  "id is required",
		},
		"email": {
			Rule: "email",
			Msg:  "invalid email",
		},
	}
	m := []struct {
		input    map[string]interface{}
		expected map[string]interface{}
		err      error
	}{
		{
			input: map[string]interface{}{
				"id": "id",
			},
			expected: map[string]interface{}{
				"id": "id",
			},
		},
		{
			input: map[string]interface{}{
				"id":    "id",
				"email": "random_string",
			},
			err: &InvalidParamsError{Errs: map[string]string{
				"email": "invalid email",
			}},
		},
		{
			input: map[string]interface{}{
				"id":        "id",
				"email":     "test@example.com",
				"extra_key": "random_value",
			},
			expected: map[string]interface{}{
				"id":    "id",
				"email": "test@example.com",
			},
		},
	}

	for _, i := range m {
		res, err := v.ValidateMap(i.input, rules)
		a.Equal(i.err, err)
		a.Equal(i.expected, res)
	}
}

func TestNewEnTranslator(t *testing.T) {
	a := assert.New(t)
	trans := NewEnTranslator()
	v := pv.New()

	if err := en.RegisterDefaultTranslations(v, trans); err != nil {
		t.Fatal(err)
	}

	err := v.Var(10, "gte=11")
	ve := err.(pv.ValidationErrors)

	a.NotNil(err)
	a.Contains(fmt.Sprintf("%v", ve.Translate(trans)), "must be 11 or greater")
}

func TestNewValidator(t *testing.T) {
	type testStruct struct {
		Desc string `validate:"is_short_desc"`
	}

	a := assert.New(t)
	m := []struct {
		s           testStruct
		tag         string
		initErr     error
		validateErr error
	}{
		{
			s:   testStruct{Desc: "1234"},
			tag: "",
			// nolint: golint
			initErr: errors.New("Function Key cannot be empty"),
		},
		{
			s:   testStruct{Desc: "1234"},
			tag: "is_short_desc",
		},
		{
			s:   testStruct{Desc: "12345"},
			tag: "is_short_desc",
			validateErr: &InvalidParamsError{Errs: map[string]string{
				"desc": "description is too long",
			}},
		},
	}

	for _, i := range m {
		v, err := New(nil, ValidationRule{
			Tag: i.tag,
			Func: func(s string) bool {
				return len(s) < 5
			},
			Msg: "description is too long",
		})

		a.Equal(i.initErr, err)
		if err == nil {
			a.Equal(i.validateErr, v.ValidateStruct(i.s))
		}
	}
}
