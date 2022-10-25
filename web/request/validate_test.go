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

func TestValidator_Register(t *testing.T) {
	err := Validate.Register("mobile", "{0}无效", func(f i18nValidator.FieldLevel) bool {
		return regexp.MustCompile(`^1[3-9][0-9]{9}$`).MatchString(f.Field().String())
	})

	if err != nil {
		t.Errorf("validate register: %v", err)
		return
	}

	v1 := &struct1{}

	if err = Validate.Struct(v1); err != nil {
		t.Errorf("validate check: %v", err)
		return
	}

	t.Log("validate registered and check succeed")
}

func TestValidator_Body(t *testing.T) {
}

func TestValidator_Struct(t *testing.T) {
}
