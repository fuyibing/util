// author: wsfuyibing <websearch@163.com>
// date: 2022-10-25

package request

import (
	i18nValidator "github.com/go-playground/validator/v10"
	"regexp"
	"testing"
)

type (
	struct1 struct {
		Mobile string `yaml:"required,min=3,max=11"`
	}

	struct2 struct {
	}
)

func ExampleValidator_Register() {
	// 1. 注册校验.
	if err := Validate.Register("mobile", "{0}无效", func(f i18nValidator.FieldLevel) bool {
		return regexp.MustCompile(`^1[3-9][0-9]{9}$`).MatchString(f.Field().String())
	}); err != nil {
		println("register:", err.Error())
		return
	}

	// 2. 定义数据.
	v := &struct {
		Mobile string `validate:"mobile" label:"手机号"`
	}{
		Mobile: "12345678901",
	}

	// 3. 校验数据.
	if err := Validate.Struct(v); err != nil {
		println("validate:", err.Error())
		return
	}

	println("complete")
}

func TestValidator_Register(t *testing.T) {
	ExampleValidator_Register()
}
