// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// Package make documents for application.
package docs

import (
	"github.com/fuyibing/util/commands/base"
)

type command struct {
	base.Command
}

func (o *command) Run(args []string) error {
	return nil
}

func New() base.CommandInterface {
	o := &command{}
	o.Initialize()
	o.SetName("docs")
	o.SetDescription("Build application documents")
	return o
}
