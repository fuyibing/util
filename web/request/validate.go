// author: wsfuyibing <websearch@163.com>
// date: 2023-02-01

package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/locales"
	i18nEN "github.com/go-playground/locales/en"
	i18nZH "github.com/go-playground/locales/zh"
	i18nTranslator "github.com/go-playground/universal-translator"
	i18nValidator "github.com/go-playground/validator/v10"
	i18nENTranslations "github.com/go-playground/validator/v10/translations/en"
	i18nZHTranslations "github.com/go-playground/validator/v10/translations/zh"

	"reflect"
)

var (
	ErrInvalidJson = fmt.Errorf("invalid json")
	Validate       *Validator
)

type (
	Validator struct {
		trans i18nTranslator.Translator
		valid *i18nValidator.Validate
	}
)

// Register
// custom tag for validation.
func (o *Validator) Register(tag, message string, check func(f i18nValidator.FieldLevel) bool) error {
	if err := o.valid.RegisterValidation(tag, check); err != nil {
		return err
	}

	return o.valid.RegisterTranslation(tag, o.trans,
		func(ut i18nTranslator.Translator) error {
			return ut.Add(tag, message, true)
		},
		func(ut i18nTranslator.Translator, fe i18nValidator.FieldError) string {
			us, ue := ut.T(fe.Tag(), fe.Field())
			if ue != nil {
				return fe.(error).Error()
			}
			return us
		},
	)
}

// Body
// unmarshal into struct and return validate result.
func (o *Validator) Body(v interface{}, body []byte) error {
	if err := json.Unmarshal(body, v); err != nil {
		return ErrInvalidJson
	}
	return o.Struct(v)
}

// Struct
// return validate result.
func (o *Validator) Struct(v interface{}) error {
	if e0 := o.valid.Struct(v); e0 != nil {
		for _, e1 := range e0.(i18nValidator.ValidationErrors) {
			return errors.New(e1.Translate(o.trans))
		}
	}
	return nil
}

// With
// translator register.
func (o *Validator) With(register func(*i18nValidator.Validate, i18nTranslator.Translator) error, fallback locales.Translator, supports ...locales.Translator) {
	if o.trans = i18nTranslator.New(fallback, supports...).GetFallback(); o.trans != nil {
		_ = register(o.valid, o.trans)
	}
}

func (o *Validator) WithEN() { o.With(i18nENTranslations.RegisterDefaultTranslations, i18nEN.New()) }
func (o *Validator) WithZH() { o.With(i18nZHTranslations.RegisterDefaultTranslations, i18nZH.New()) }

func (o *Validator) init() *Validator {
	o.valid = i18nValidator.New()

	o.valid.RegisterTagNameFunc(func(field reflect.StructField) string {
		if v := field.Tag.Get("label"); v != "" {
			return v
		}
		return field.Name
	})

	o.WithZH()
	return o
}
