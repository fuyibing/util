// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// Package make documents for application.
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
	o.SetDescription("Build application documents")
	return o
}

func (o *command) Run(args []string) error {
	if err := o.ParseArguments(args); err != nil {
		return err
	}
	return errors.New(fmt.Sprintf("%s: todo Run() method", o.Name()))
}
