// author: wsfuyibing <websearch@163.com>
// date: 2021-02-16

package base

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	OptionRegexpDouble = regexp.MustCompile(`^--([a-zA-Z0-9][_a-zA-Z0-9\-]*)[=]?(.*)`)
	OptionRegexpSingle = regexp.MustCompile(`^-([_a-zA-Z0-9]+)$`)
)

// Option mode list.
const (
	OptionModeOptional OptionMode = iota
	OptionModeRequired
)

// Option mode type.
type OptionMode int

// Value mode list.
const (
	OptionValueModeNone OptionValueMode = iota
	OptionValueModeString
	OptionValueModeInteger
)

// Value mode type.
type OptionValueMode int

// Option interface.
type OptionInterface interface {
	IsOptional() bool
	IsRequired() bool
	Name() string
	ShortName() string
	SetDefaultValue(interface{}) OptionInterface
	SetDescription(string) OptionInterface
	SetShortName(string) OptionInterface
	SetValue(interface{}) OptionInterface
	ToBool() (val bool, err error)
	ToInt() (val int, err error)
	ToInt64() (val int64, err error)
	ToString() (val string, err error)
	Usage(string)
}

// Option struct.
type option struct {
	mode         OptionMode
	valueMode    OptionValueMode
	name         string
	value        interface{}
	shortName    string
	defaultValue interface{}
	description  string
}

// Create option instance.
func NewOption(name string, mode OptionMode, valueMode OptionValueMode) OptionInterface {
	return &option{
		name:      name,
		mode:      mode,
		valueMode: valueMode,
	}
}

// Is optional mode.
func (o *option) IsOptional() bool {
	return o.mode == OptionModeOptional
}

// Is requirement mode.
func (o *option) IsRequired() bool {
	return o.mode == OptionModeRequired
}

// Return option name.
func (o *option) Name() string {
	return o.name
}

// Return option short name.
func (o *option) ShortName() string {
	return o.shortName
}

// Set default value.
func (o *option) SetDefaultValue(defaultValue interface{}) OptionInterface {
	o.defaultValue = defaultValue
	return o
}

// Set description.
func (o *option) SetDescription(description string) OptionInterface {
	o.description = description
	return o
}

// Set option short name.
func (o *option) SetShortName(shortName string) OptionInterface {
	o.shortName = shortName
	return o
}

// Set option value.
func (o *option) SetValue(value interface{}) OptionInterface {
	o.value = value
	return o
}

// To boolean value.
func (o *option) ToBool() (bool, error) {
	// can not convert to boolean.
	if o.valueMode != OptionValueModeNone {
		return false,
			errors.New(fmt.Sprintf("Option %s: can not convert not boolean option to boolean", o.name))
	}
	// return false if not specified.
	if o.value == nil {
		return false, nil
	}
	// read option value.
	v := fmt.Sprintf("%v", o.value)
	// return true if value is empty.
	if v == "" {
		return true, nil
	}
	// parse string to boolean.
	b, err := strconv.ParseBool(strings.ToLower(fmt.Sprintf("%v", o.value)))
	if err != nil {
		return false,
			errors.New(fmt.Sprintf("Option %s: can not convert to boolean: %s", o.name, v))
	}
	return b, nil
}

// To integer value.
func (o *option) ToInt() (int, error) {
	// not integer.
	if o.valueMode != OptionValueModeInteger {
		return 0,
			errors.New(fmt.Sprintf("Option %s: can not convert not integer option to integer", o.name))
	}
	// assign value.
	var v interface{}
	if o.value != nil {
		v = o.value
	} else if o.defaultValue != nil {
		v = o.defaultValue
	} else {
		return 0,
			errors.New(fmt.Sprintf("Option %s: integer value not specified", o.name))
	}
	// parse to integer
	vi, err := strconv.ParseInt(fmt.Sprintf("%v", v), 0, 32)
	if err != nil {
		return 0,
			errors.New(fmt.Sprintf("Option %s: can not convert to integer: %v", o.name, v))
	}
	return int(vi), nil
}

// To integer 64 value.
func (o *option) ToInt64() (int64, error) {
	if o.valueMode != OptionValueModeInteger {
		return 0,
			errors.New(fmt.Sprintf("Option %s: can not convert not integer option to integer", o.name))
	}
	// assign value.
	var v interface{}
	if o.value != nil {
		v = o.value
	} else if o.defaultValue != nil {
		v = o.defaultValue
	} else {
		return 0,
			errors.New(fmt.Sprintf("Option %s: integer value not specified", o.name))
	}
	// parse to integer
	vi, err := strconv.ParseInt(fmt.Sprintf("%v", v), 0, 64)
	if err != nil {
		return 0,
			errors.New(fmt.Sprintf("Option %s: can not convert to integer: %v", o.name, v))
	}
	return vi, nil
}

// To string value.
func (o *option) ToString() (string, error) {
	if o.valueMode != OptionValueModeString {
		return "",
			errors.New(fmt.Sprintf("Option %s: can not convert not string option to string", o.name))
	}
	var v interface{}
	if o.value != nil {
		v = o.value
	} else if o.defaultValue != nil {
		v = o.defaultValue
	} else {
		if o.IsRequired() {
			return "",
				errors.New(fmt.Sprintf("Option %s: string value not specified", o.name))
		} else {
			return "", nil
		}
	}
	return fmt.Sprintf("%v", v), nil
}

// Print usage.
func (o *option) Usage(prefix string) {
	var s = prefix
	// 1. short option name
	if o.shortName != "" {
		s = fmt.Sprintf("%s -%s,", prefix, o.shortName)
	} else {
		s = fmt.Sprintf("%s    ", prefix)
	}
	// 2. option name
	s += fmt.Sprintf("--%s", o.name)
	v := ""
	if o.defaultValue != nil {
		v = fmt.Sprintf("%v", o.defaultValue)
	}
	// 3. option value.
	if o.valueMode != OptionValueModeNone {
		if o.mode == OptionModeRequired {
			// 3.1 required
			if v == "" {
				if o.valueMode == OptionValueModeString {
					v = "STR"
				} else if o.valueMode == OptionValueModeInteger {
					v = "INT"
				} else {
					v = "VAL"
				}
			}
			s += fmt.Sprintf("=<%s>", v)
		} else if o.mode == OptionModeOptional {
			// 3.2 optional
			if v == "" {
				if o.valueMode == OptionValueModeString {
					v = "STR"
				} else if o.valueMode == OptionValueModeInteger {
					v = "INT"
				} else {
					v = "VAL"
				}
			}
			s += fmt.Sprintf("[=%s]", v)
		}
	}
	// n. print usage.
	fmt.Printf("%-48s %s\n", s, o.description)
}
