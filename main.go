// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package main

import (
	"fmt"

	"github.com/fuyibing/util/commands"
)

func main() {
	m := commands.Default()
	if err := m.Run(); err != nil {
		fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, 0, 33, 41, err.Error(), 0x1B)
		return
	}
}
