// author: wsfuyibing <websearch@163.com>
// date: 2021-02-14

package main

import (
	"fmt"

	"github.com/fuyibing/util/commands"
)

func main() {
	m := commands.Default()
	if err := m.Run(nil); err != nil {
		fmt.Printf("%v.\n", err)
		return
	}
}
