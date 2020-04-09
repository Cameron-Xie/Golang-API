package validator

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	pv "github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	ValidateStruct(interface{}) error
	ValidateMap(i map[string]interface{}, rules map[string]MapRule) (map[string]interface{}, error)
}

type ValidationRuleFunc func(string) bool

type ValidationRule struct {
	Tag  string
	Func ValidationRuleFunc
	Msg  string
}

type MapRule struct {
	Rule string
	Msg  string
}

type validator struct {
	v *pv.Validate
	t ut.Translator
}

func (v *validator) ValidateStruct(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		l := make(map[string]string)

		for _, e := range err.(pv.ValidationErrors) {
			l[strings.ToLower(e.Field())] = strings.ToLower(e.Translate(v.t))
		}

		return &InvalidParamsError{Errs: l}
	}

	return nil
}

func (v *validator) ValidateMap(i map[string]interface{}, rules map[string]MapRule) (map[string]interface{}, error) {
	res, errs := make(map[string]interface{}), make(map[string]string)
	for key, val := range i {
		if _, ok := rules[key]; !ok {
			continue
		}

		r := rules[key]
		if err := v.v.Var(val, r.Rule); err != nil {
			errs[key] = r.Msg
			continue
		}

		res[key] = val
	}

	if len(errs) == 0 {
		return res, nil
	}

	return nil, &InvalidParamsError{Errs: errs}
}

func NewEnTranslator() ut.Translator {
	enl := en.New()
	uni := ut.New(enl, enl)
	trans, _ := uni.GetTranslator("en")

	return trans
}

func New(trans ut.Translator, rules ...ValidationRule) (Validator, error) {
	v := pv.New()

	if trans == nil {
		trans = NewEnTranslator()
	}
	_ = enTrans.RegisterDefaultTranslations(v, trans)

	for _, r := range rules {
		err := v.RegisterValidation(r.Tag, createValidationRule(r.Func))
		if err != nil {
			return nil, err
		}

		if trans != nil {
			_ = v.RegisterTranslation(r.Tag, trans, registerFn(r.Tag, r.Msg), translationFn)
		}
	}

	return &validator{
		v: v,
		t: trans,
	}, nil
}

func createValidationRule(fn ValidationRuleFunc) pv.Func {
	return func(fl pv.FieldLevel) bool {
		return fn(fl.Field().String())
	}
}

func registerFn(tag, translation string) pv.RegisterTranslationsFunc {
	return func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}
}

func translationFn(tran ut.Translator, fe pv.FieldError) string {
	t, _ := tran.T(fe.Tag(), fe.Field())

	return t
}
