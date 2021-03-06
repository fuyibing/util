// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 创建项目文档.
package docs

import (
	"errors"
	"fmt"

	"github.com/fuyibing/util/cmds/base"
)

type command struct {
	base.Command
}

func New() base.CommandInterface {
	o := &command{}
	o.Initialize()
	o.SetName("docs")
	o.SetDescription("Create application documents")
	// append options
	o.AddOption(
		base.NewOption("path", base.OptionModeOptional, base.OptionValueModeString).
			SetShortName("p").
			SetDefaultValue("./framework").
			SetDescription("Source file directory (default: ./framework)"),
		base.NewOption("target", base.OptionModeOptional, base.OptionValueModeString).
			SetShortName("t").
			SetDefaultValue("./docs").
			SetDescription("Created file save to (default: ./docs)"),
	)
	return o
}

// Run command.
func (o *command) Run(manager base.ManagerInterface, args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("Command %s: TODO", o.GetName()))
}
