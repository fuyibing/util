// author: wsfuyibing <websearch@163.com>
// date: 2021-02-15

// 命令: 命令管理器.
package cmds

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/fuyibing/util/cmds/base"
	"github.com/fuyibing/util/cmds/help"
	"github.com/fuyibing/util/cmds/kv"
	"github.com/fuyibing/util/cmds/makes"
)

const (
	DefaultCommandName    = "fyb"
	DefaultCommandVersion = "0.0"
)

// 管理器结构体.
type management struct {
	commands map[string]base.CommandInterface
	mu       *sync.RWMutex
	name     string
	version  string
}

// 添加命令.
func (o *management) AddCommand(cs ...base.CommandInterface) base.ManagerInterface {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, c := range cs {
		o.commands[c.GetName()] = c
	}
	return o
}

// 读取命令.
func (o *management) GetCommand(name string) base.CommandInterface {
	o.mu.RLock()
	defer o.mu.RUnlock()
	if cmd, ok := o.commands[name]; ok {
		return cmd
	}
	return nil
}

// 读取命令列表.
func (o *management) GetCommands() map[string]base.CommandInterface {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.commands
}

// 读取项目名称.
func (o *management) GetName() string { return o.name }

// 读取项目版本号.
func (o *management) GetVersion() string { return o.version }

// 运行指定命令.
func (o *management) Run(args ...string) error {
	// 1. initialize arguments.
	if args == nil {
		args = os.Args
	}
	// 2. arguments length less than 2 fields.
	if args == nil || len(args) < 2 {
		return errors.New(fmt.Sprintf("Command: command name not specified"))
	}
	// 3. command name.
	name := args[1]
	// 3.1 run added command.
	o.mu.RLock()
	defer o.mu.RUnlock()
	if c, ok := o.commands[name]; ok {
		args[0] = "go run main.go"
		return c.Run(o, args)
	}
	// 4. return error if not added.
	return errors.New(fmt.Sprintf("Command: command name not defined: %s", name))
}

// 初始化配置参数.
func (o *management) initialize() {
	o.mu = new(sync.RWMutex)
	o.commands = make(map[string]base.CommandInterface)
	o.name = DefaultCommandName
	o.version = DefaultCommandVersion
	// parse yaml.
	var tmp = &struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}{}
	for _, file := range []string{"./tmp/framework.yaml", "../tmp/framework.yaml", "./config/framework.yaml", "../config/framework.yaml"} {
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		if err = yaml.Unmarshal(bs, tmp); err != nil {
			continue
		}
		break
	}
	if tmp.Name != "" {
		o.name = tmp.Name
	}
	if tmp.Version != "" {
		o.version = tmp.Version
	}
}

// 创建默认管理器.
func Default() base.ManagerInterface {
	o := New()
	o.AddCommand(
		makes.New(),
		kv.New(),
		help.New(),
	)
	return o
}

// 创建管理器.
func New() base.ManagerInterface {
	o := &management{}
	o.initialize()
	return o
}
