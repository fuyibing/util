// author: wsfuyibing <websearch@163.com>
// date: 2021-02-15

// Package command line manager.
package commands

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/fuyibing/util/commands/base"
	"github.com/fuyibing/util/commands/docs"
	"github.com/fuyibing/util/commands/makes"
)

// Command line manager interface.
type ManagerInterface interface {
	AddCommand(...base.CommandInterface) ManagerInterface
	Run([]string) error
}

// 命令行管理器结构体.
type management struct {
	mu       *sync.RWMutex
	commands map[string]base.CommandInterface
}

// 添加命令接口.
func (o *management) AddCommand(cs ...base.CommandInterface) ManagerInterface {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, c := range cs {
		o.commands[c.Name()] = c
	}
	return o
}

// 运行命令行工具.
func (o *management) Run(args []string) error {
	// 1. 初始化入参.
	if args == nil {
		args = os.Args
	}
	// 2. 命令参数不少于2位.
	if args == nil || len(args) < 2 {
		return errors.New(fmt.Sprintf("Command: name not specified"))
	}
	// 3. 命令名称.
	name := args[1]
	// 3.1 今天已注册.
	o.mu.RLock()
	defer o.mu.RUnlock()
	if c, ok := o.commands[name]; ok {
		return c.Run(args)
	}
	// 4. 返回未注册错误
	return errors.New(fmt.Sprintf("Command: undefined command: %s", name))
}

// 创建默认命令行管理器实例.
func Default() ManagerInterface {
	o := New()
	o.AddCommand(makes.New(), docs.New())
	return o
}

// 创建命令行管理器实例.
func New() ManagerInterface {
	o := &management{}
	o.mu = new(sync.RWMutex)
	o.commands = make(map[string]base.CommandInterface)
	return o
}
