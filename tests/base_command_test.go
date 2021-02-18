// author: wsfuyibing <websearch@163.com>
// date: 2021-02-15

package tests

import (
	"testing"

	"github.com/fuyibing/util/commands/base"
	"github.com/fuyibing/util/commands/makes"
)

func TestBaseCommandMakes(t *testing.T) {
	c := makes.New()
	if err := c.Run([]string{
		"xmbs",
		"make",
		"--type=model",
		"--name=example",
		"--override",
	}); err != nil {
		t.Errorf("Command error: %v.", err)
		return
	}

	if opt, ok := c.GetOption("override"); ok {
		v, err := opt.ToBool()
		if err != nil {
			t.Errorf("Option error: %v.", err)
			return
		}
		t.Logf("Option bool: %v.", v)
	}


	t.Logf("Completed.")
}

func TestBaseCommandConsul(t *testing.T) {
	c := base.NewCommand("kv")
	c.AddOption(
		base.NewOption("app", base.OptionModeRequired, base.OptionValueModeString).
			SetShortName("a").
			SetDescription("Specify your application name"),
		base.NewOption("consul", base.OptionModeRequired, base.OptionValueModeString).
			SetShortName("c").
			SetDescription("Consul address"),
	)
	c.Usage()
}

// import (
// 	"testing"
//
// 	"github.com/fuyibing/util/commands/base2"
// 	"github.com/fuyibing/util/commands/makes"
// )
//
// func TestCommandManager(t *testing.T) {
// 	ok, err := base2.Manager.AddCommand(makes.New()).Run(
// 		"app",
// 		"make",
// 		"-t",
// 		"model",
// 	)
//
// 	if ok {
// 		if err != nil {
// 			t.Errorf("command: %v.", err)
// 			return
// 		}
// 		t.Logf("command completed.")
// 		return
// 	}
//
// 	t.Logf("unknown command")
// }
//
// func TestCommandMake(t *testing.T) {
// 	cmd := makes.New()
//
// 	if err := cmd.Run(nil); err != nil {
// 		t.Errorf("run {%s} error: %v.", cmd.GetName(), err)
// 		return
// 	}
// 	t.Logf("run {%s} completed.", cmd.GetName())
// 	cmd.Usage()
//
// }
//
// // import (
// // 	"testing"
// //
// // 	"github.com/fuyibing/util/commands/base2"
// // )
// //
// // type MyCommand struct {
// // 	base2.Command
// // }
// //
// // func NewCommand(name string) base2.Command {
// // 	o := &MyCommand{}
// // 	return o
// // }
// //
// // func TestBaseCommamnd1(t *testing.T) {
// // 	t.Logf("---- ---- [ command ] ---- ----")
// //
// // 	cmd := NewCommand("test")
// // 	cmd.AddDefinition(base2.NewDefinition("host"))
// // 	cmd.Run(nil)
// //
// // }
