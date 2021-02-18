// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 创建项目文档.
package docs

import (
	"errors"
	"fmt"

	"github.com/fuyibing/util/commands/base"
)

type command struct {
	base.Command
}

func New() base.CommandInterface {
	o := &command{}
	o.Initialize()
	o.SetName("docs")
	o.SetDescription("创建项目文档")
	return o
}

func (o *command) Run(manager base.ManagerInterface, args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("Command %s: TODO", o.GetName()))
}
