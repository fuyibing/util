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
	Validate() error
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
func (o *option) IsOptional() bool { return o.mode == OptionModeOptional }

// Is requirement mode.
func (o *option) IsRequired() bool { return o.mode == OptionModeRequired }

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
	// only work on none value.
	if o.valueMode != OptionValueModeNone {
		return false, errors.New("not none value option")
	}
	// not specified.
	if o.value == nil {
		return false, nil
	}
	// return default if option value is empty string.
	v := fmt.Sprintf("%v", o.value)
	if v == "" {
		return true, nil
	}
	// parse string to boolean.
	b, err := strconv.ParseBool(strings.ToLower(fmt.Sprintf("%v", o.value)))
	if err != nil {
		return false, errors.New("not boolean value")
	}
	return b, nil
}

// To integer value.
func (o *option) ToInt() (int, error) {
	// not integer.
	if o.valueMode != OptionValueModeInteger {
		return 0, errors.New("not integer value option")
	}
	// assign value.
	var v interface{}
	if o.value != nil {
		v = o.value
	} else if o.defaultValue != nil {
		v = o.defaultValue
	} else {
		return 0, errors.New("integer value not specified")
	}
	// parse to integer
	vi, err := strconv.ParseInt(fmt.Sprintf("%v", v), 0, 32)
	if err != nil {
		return 0, err
	}
	return int(vi), nil
}

// To integer 64 value.
func (o *option) ToInt64() (int64, error) {
	if o.valueMode != OptionValueModeInteger {
		return 0, errors.New("not integer value option")
	}
	// assign value.
	var v interface{}
	if o.value != nil {
		v = o.value
	} else if o.defaultValue != nil {
		v = o.defaultValue
	} else {
		return 0, errors.New("integer value not specified")
	}
	// parse to integer
	vi, err := strconv.ParseInt(fmt.Sprintf("%v", v), 0, 64)
	if err != nil {
		return 0, err
	}
	return vi, nil
}

// To string value.
func (o *option) ToString() (string, error) {
	if o.valueMode != OptionValueModeString {
		return "", errors.New("not string value option")
	}
	var v interface{}
	if o.value != nil {
		v = o.value
	} else if o.defaultValue != nil {
		v = o.defaultValue
	} else {
		return "", errors.New("string value not specified")
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

// Validate option.
func (o *option) Validate() error {
	if o.mode == OptionModeRequired {
		// 1.1 none value mode.
		if o.valueMode == OptionValueModeNone {
			return nil
		}
		// 1.2 string
		if o.value == nil {
			return errors.New(fmt.Sprintf("option %s not speicified", o.name))
		}
	}
	return nil
}
