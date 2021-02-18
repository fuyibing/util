// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 基础依赖.
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
	GetDescription() string
	GetName() string
	GetOption(string) (OptionInterface, error)
	Info(string, ...interface{})
	Initialize() CommandInterface
	IsHidden() bool
	ParseArguments([]string) error
	Run(ManagerInterface, []string) error
	SetDescription(string) CommandInterface
	SetHidden(bool) CommandInterface
	SetName(string) CommandInterface
	Usage(ManagerInterface)
}

// Command struct.
type Command struct {
	hidden      bool
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

// Return command description.
func (o *Command) GetDescription() string {
	return o.description
}

// Return command name.
func (o *Command) GetName() string {
	return o.name
}

// Get option by name.
func (o *Command) GetOption(name string) (OptionInterface, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if opt, ok := o.options[name]; ok {
		return opt, nil
	}
	return nil, errors.New(fmt.Sprintf("Command %s: option not defined: %s", o.name, name))
}

// Print info.
func (o *Command) Info(text string, args ...interface{}) {
	println(fmt.Sprintf(text, args...))
}

// Initialize command fields.
func (o *Command) Initialize() CommandInterface {
	o.mu = new(sync.RWMutex)
	o.options = make(map[string]OptionInterface)
	return o
}

// Command is hidden.
func (o *Command) IsHidden() bool {
	return o.hidden
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
			return errors.New(fmt.Sprintf("Command %s: option not specified: %s", o.name, opt.Name()))
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
func (o *Command) Run(manager ManagerInterface, args []string) error {
	return errors.New(fmt.Sprintf("Command %s: Run() method override", o.name))
}

// Set command description.
func (o *Command) SetDescription(description string) CommandInterface {
	o.description = description
	return o
}

// Set status hidden.
func (o *Command) SetHidden(hidden bool) CommandInterface {
	o.hidden = hidden
	return o
}

// Set command name.
func (o *Command) SetName(name string) CommandInterface {
	o.name = name
	return o
}

// Print usage.
func (o *Command) Usage(manager ManagerInterface) {
	// 1. print usage.
	fmt.Printf("Application : %s/%s\n", manager.GetName(), manager.GetVersion())
	fmt.Printf("Usage       : %s %s [OPTOINS]\n", "go run main.go", o.name)
	// 2. print options
	o.mu.RLock()
	defer o.mu.RUnlock()
	if len(o.options) > 0 {
		i := 0
		for _, c := range o.options {
			if i++; i == 1 {
				c.Usage("Options     :")
			} else {
				c.Usage("            :")
			}
		}
	}
}
