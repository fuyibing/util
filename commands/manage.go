// author: wsfuyibing <websearch@163.com>
// date: 2021-02-15

// 命令: 命令管理器.
package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/fuyibing/util/commands/base"
	"github.com/fuyibing/util/commands/docs"
	"github.com/fuyibing/util/commands/help"
	"github.com/fuyibing/util/commands/kv"
	"github.com/fuyibing/util/commands/makes"
)

const (
	DefaultCommandName    = "fyb"
	DefaultCommandVersion = "0.0"
)

// Command line manager struct.
type management struct {
	commands map[string]base.CommandInterface
	mu       *sync.RWMutex
	name     string
	version  string
}

// Add command to manager.
func (o *management) AddCommand(cs ...base.CommandInterface) base.ManagerInterface {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, c := range cs {
		o.commands[c.GetName()] = c
	}
	return o
}

// Get added command.
func (o *management) GetCommand(name string) base.CommandInterface {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if cmd, ok := o.commands[name]; ok {
		return cmd
	}
	return nil
}

// Get added commands.
func (o *management) GetCommands() map[string]base.CommandInterface {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.commands
}

// Get application name.
func (o *management) GetName() string {
	return o.name
}

// Get application version.
func (o *management) GetVersion() string {
	return o.version
}

// Run command manager.
func (o *management) Run(args ...string) error {
	// 1. initialize arguments.
	if args == nil {
		args = os.Args
	}
	// 2. arguments length less than 2 fields.
	if args == nil || len(args) < 2 {
		return errors.New(fmt.Sprintf("Command: no command name"))
	}
	// 3. command name.
	name := args[1]
	// 3.1 run added command.
	o.mu.RLock()
	defer o.mu.RUnlock()
	if c, ok := o.commands[name]; ok {
		return c.Run(o, args)
	}
	// 4. return error if not added.
	return errors.New(fmt.Sprintf("Command: undefined command: %s", name))
}

// Initialize manager instance.
func (o *management) initialize() {
	// reset name and version.
	o.mu = new(sync.RWMutex)
	o.commands = make(map[string]base.CommandInterface)
	o.name = DefaultCommandName
	o.version = DefaultCommandVersion
	// parse yaml.
	tmp := &struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}{}
	for _, file := range []string{"./config/app.yaml", "../config/app.yaml"} {
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		if err = yaml.Unmarshal(bs, tmp); err != nil {
			continue
		}
	}
	if tmp.Name != "" {
		o.name = tmp.Name
	}
	if tmp.Version != "" {
		o.version = tmp.Version
	}
}

// Create default manager.
func Default() base.ManagerInterface {
	o := New()
	o.AddCommand(makes.New(), docs.New(), kv.New(), help.New())
	return o
}

// Create empty manager.
func New() base.ManagerInterface {
	o := &management{}
	o.initialize()
	return o
}
