// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

package base

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Command interface.
type CommandInterface interface {
	AddOption(...OptionInterface) CommandInterface
	GetOption(string) (OptionInterface, bool)
	Initialize() CommandInterface
	Description() string
	Name() string
	ParseArguments([]string) error
	Run([]string) error
	SetDescription(string) CommandInterface
	SetName(string) CommandInterface
	Usage()
}

// Command struct.
type Command struct {
	mu          *sync.RWMutex
	name        string
	description string
	options     map[string]OptionInterface
}

// Create command instance.
func NewCommand(name string) CommandInterface {
	o := &Command{}
	o.SetName(name)
	o.Initialize()
	return o
}

// Add option.
func (o *Command) AddOption(opts ...OptionInterface) CommandInterface {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, x := range opts {
		o.options[x.Name()] = x
	}
	return o
}

// Get option by name.
func (o *Command) GetOption(name string) (OptionInterface, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if opt, ok := o.options[name]; ok {
		return opt, true
	}
	return nil, false
}

// Initialize command fields.
func (o *Command) Initialize() CommandInterface {
	o.mu = new(sync.RWMutex)
	o.options = make(map[string]OptionInterface)
	return o
}

// Return command description.
func (o *Command) Description() string {
	return o.description
}

// Return command name.
func (o *Command) Name() string {
	return o.name
}

// Parse arguments.
func (o *Command) ParseArguments(args []string) error {
	// 1. lock and release.
	o.mu.RLock()
	defer o.mu.RUnlock()
	// 2. variables for options loop.
	var found = false
	var maxIndex = len(args) - 1
	var value = ""
	// 3. loop options.
	for _, opt := range o.options {
		// parse option.
		for index, arg := range args {
			// match: name.
			if m := OptionRegexpDouble.FindStringSubmatch(arg); len(m) == 3 {
				if m[1] = strings.TrimSpace(m[1]); m[1] == opt.Name() {
					m[2] = strings.TrimSpace(m[2])
					found = true
					value = m[2]
					break
				}
			}
			// match: short name.
			if shortName := opt.ShortName(); shortName != "" {
				if m := OptionRegexpSingle.FindStringSubmatch(arg); len(m) == 2 {
					for mi := 0; mi < len(m[1]); mi++ {
						if ms := string(m[1][mi]); ms == shortName {
							found = true
							if index < maxIndex {
								value = args[index+1]
							}
							break
						}
					}
				}
			}
		}
		// return error if required option not specified.
		if opt.IsRequired() && !found {
			return errors.New(fmt.Sprintf("%s: option %s not specified", o.name, opt.Name()))
		}
		// assign option value.
		if found {
			found = false
			opt.SetValue(value)
		}
		// clean current value.
		if value != "" {
			value = ""
		}
	}

	return nil
}

// Run command.
func (o *Command) Run(args []string) error {
	return errors.New(fmt.Sprintf("%s: Run() method not defined", o.name))
}

func (o *Command) SetDescription(description string) CommandInterface {
	o.description = description
	return o
}

// Set command name.
func (o *Command) SetName(name string) CommandInterface {
	o.name = name
	return o
}

// Print usage.
func (o *Command) Usage() {
	// 1. print usage.
	fmt.Printf("Usage: %s %s [OPTIOINS]\n", "go run main.go", o.name)
	// 2. print options
	o.mu.RLock()
	defer o.mu.RUnlock()
	if len(o.options) > 0 {
		fmt.Print("Options: \n")
		for _, c := range o.options {
			c.Usage()
		}
	}
}
