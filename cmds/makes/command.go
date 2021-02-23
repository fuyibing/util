// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 脚手架.
package makes

import (
	"github.com/fuyibing/util/cmds/base"
)

type command struct {
	base.Command
}

func New() base.CommandInterface {
	o := new(command)
	o.Initialize()
	o.SetName("make")
	o.SetDescription("Create application file")
	// 2. add option.
	o.AddOption(
		base.NewOption("type", base.OptionModeRequired, base.OptionValueModeString).
			SetShortName("t").
			SetDescription("Specify your file type (accept: model|service|logic|controller|path)"),
		base.NewOption("name", base.OptionModeRequired, base.OptionValueModeString).
			SetShortName("n").
			SetDescription("Specify your file name"),
		base.NewOption("table-name", base.OptionModeOptional, base.OptionValueModeString).
			SetDescription("Specify table name for make model"),
		base.NewOption("path", base.OptionModeOptional, base.OptionValueModeString).
			SetShortName("p").
			SetDefaultValue("./framework").
			SetDescription("Specify your root path (default: ./framework)"),
		base.NewOption("override", base.OptionModeOptional, base.OptionValueModeNone).
			SetDescription("Override if file exists"),
		base.NewOption("list", base.OptionModeOptional, base.OptionValueModeNone).
			SetDescription("Print list"),
	)
	return o
}

// Parse arguments.
func (o *command) Run(manager base.ManagerInterface, args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	mng := &management{args:args}
	mng.initialize(o)
	return mng.run()
}
