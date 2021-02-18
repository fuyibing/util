// author: wsfuyibing <websearch@163.com>
// date: 2021-02-15

package tests
//
// import (
// 	"testing"
//
// 	"github.com/fuyibing/util/commands/base2"
// )
//
// var ds map[string]base2.DefinitionInterface
//
// func TestBaseDefinition1(t *testing.T) {
// 	ds = make(map[string]base2.DefinitionInterface)
// 	t.Logf("---- ---- [ definition ] ---- ----")
//
// 	add(
// 		base2.NewDefinition("host").
// 			SetShortName("h").
// 			SetDefault("127.0.0.1").
// 			SetRequired(true).
// 			SetDescription("server host addr, IPv4"),
// 	)
// 	add(
// 		base2.NewDefinition("bool").
// 			SetDescription("boolean option"),
// 	)
// 	add(
// 		base2.NewDefinition("override").SetShortName("o").AsNil().
// 			SetDescription("override if exist"),
// 	)
//
// 	usage()
// }
//
// func add(d base2.DefinitionInterface) {
// 	ds[d.Name()] = d
// }
//
// func usage() {
// 	for _, d := range ds {
// 		d.Usage()
// 	}
// }
