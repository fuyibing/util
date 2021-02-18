// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// Package make files for application.
//
// Dependent on iris framework, allow create model、
// service、logic、controller.
package makes

import (
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
	o.SetDescription("Build application files for iris framework.")
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
func (o *command) Run(args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	return nil
}
