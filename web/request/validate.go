// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package request

import (
	"encoding/json"
	"errors"
	"fmt"
	i18n "github.com/go-playground/locales/zh"
	i18nTranslator "github.com/go-playground/universal-translator"
	i18nValidator "github.com/go-playground/validator/v10"
	i18nTranslations "github.com/go-playground/validator/v10/translations/zh"

	"reflect"
)

var (
	// Validate
	// 校验实例.
	Validate *Validator

	errInvalidJson = fmt.Errorf("无效JSON入参")
)

type (
	// Validator
	// 校验结构体.
	Validator struct {
		trans i18nTranslator.Translator
		valid *i18nValidator.Validate
	}
)

// Register
// 注册校验.
func (o *Validator) Register(tag, message string, check func(f i18nValidator.FieldLevel) bool) (err error) {
	if err = o.valid.RegisterValidation(tag, check); err == nil {
		err = o.valid.RegisterTranslation(tag, o.trans,
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
	return
}

// Body
// 校验入参.
func (o *Validator) Body(v interface{}, body []byte) error {
	if err := json.Unmarshal(body, v); err != nil {
		return errInvalidJson
	}
	return o.Struct(v)
}

// Struct
// 校验结构体.
func (o *Validator) Struct(v interface{}) error {
	if e0 := o.valid.Struct(v); e0 != nil {
		for _, e1 := range e0.(i18nValidator.ValidationErrors) {
			return errors.New(e1.Translate(o.trans))
		}
	}
	return nil
}

// 构造实例.
func (o *Validator) init() *Validator {
	// 1. 创建实例.
	o.valid = i18nValidator.New()

	// 2. 解析标签.
	//
	//   type Example struct{
	//       Mobile string `label:"手机号"`
	//   }
	o.valid.RegisterTagNameFunc(func(field reflect.StructField) string {
		if v := field.Tag.Get("label"); v != "" {
			return v
		}
		return field.Name
	})

	// 2. 绑定中文.
	var found bool
	if o.trans, found = i18nTranslator.New(i18n.New(), i18n.New()).GetTranslator("zh"); found {
		_ = i18nTranslations.RegisterDefaultTranslations(o.valid, o.trans)
	}

	// n. 完成创建.
	return o
}
