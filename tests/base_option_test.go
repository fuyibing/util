// author: wsfuyibing <websearch@163.com>
// date: 2021-02-16

package tests

import (
	"testing"

	"github.com/fuyibing/util/cmds/base"
)

func TestBaseOption(t *testing.T) {
	base.NewOption("host", base.OptionModeRequired, base.OptionValueModeString).
		SetShortName("h").
		SetDescription("server address, IPv4").
		Usage()
	base.NewOption("port", base.OptionModeRequired, base.OptionValueModeInteger).
		SetShortName("p").
		SetDescription("server listen port").
		Usage()
	base.NewOption("delay", base.OptionModeOptional, base.OptionValueModeInteger).
		SetDefaultValue(60).
		Usage()
	base.NewOption("force", base.OptionModeOptional, base.OptionValueModeNone).
		SetShortName("f").
		Usage()
}
