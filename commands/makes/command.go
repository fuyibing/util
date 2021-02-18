// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 脚手架.
package makes

import (
	"errors"
	"fmt"

	"github.com/fuyibing/util/commands/base"
)

type command struct {
	base.Command
}

// Create MAKE command.
func New() base.CommandInterface {
	// create empty instance.
	// initialize fields and set command name.
	o := new(command)
	o.Initialize()
	o.SetName("make")
	o.SetDescription("脚手架, 创建Model、Service、Logic、Controller文件")
	// 2. add option.
	o.AddOption(
		base.NewOption("type", base.OptionModeRequired, base.OptionValueModeString).
			SetDescription("Specify your file type, accept: model|service|logic|controller"),
		base.NewOption("name", base.OptionModeRequired, base.OptionValueModeString).
			SetDescription("Specify your file name"),
		base.NewOption("override", base.OptionModeOptional, base.OptionValueModeNone).
			SetDescription("Override if file exists"),
	)
	return o
}

// Parse arguments.
func (o *command) Run(manager base.ManagerInterface, args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("Command %s: TODO", o.GetName()))
}
