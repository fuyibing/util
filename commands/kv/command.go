// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// Package make files for application.
//
// Dependent on iris framework, allow create model、
// service、logic、controller.
package kv

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
	o.SetName("kv")
	o.SetDescription("Build application config use consul kv.")
	return o
}

// Parse arguments.
func (o *command) Run(manager base.ManagerInterface, args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("%s: todo Run() method", o.GetName()))
}
