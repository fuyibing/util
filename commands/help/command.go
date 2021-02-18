// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 帮助向导.
package help

import (
	"fmt"

	"github.com/fuyibing/util/commands/base"
)

type command struct {
	base.Command
}

func New() base.CommandInterface {
	o := new(command)
	o.Initialize()
	o.SetName("help")
	o.SetHidden(true)
	return o
}

func (o *command) Run(manager base.ManagerInterface, args []string) error {
	// Help command.
	if len(args) >= 3 {
		key := args[2]
		if cmd := manager.GetCommand(key); cmd != nil {
			cmd.Usage(manager)
			return nil
		}
	}
	// Command list.
	fmt.Printf("应用    : %s/%s\n", manager.GetName(), manager.GetVersion())
	fmt.Printf("用法    : go run main.go <命令> [选项]\n")
	i := 0
	for _, c := range manager.GetCommands() {
		if c.IsHidden() {
			continue
		}
		if i++; i == 1 {
			fmt.Printf("命令 %2d : %-34s %s\n", i, c.GetName(), c.GetDescription())
		} else {
			fmt.Printf("     %2d : %-34s %s\n", i, c.GetName(), c.GetDescription())
		}
	}
	return nil
}
