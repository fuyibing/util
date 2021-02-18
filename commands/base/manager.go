// author: wsfuyibing <websearch@163.com>
// date: 2021-02-18

package base

// Command manager interface.
type ManagerInterface interface {
	AddCommand(...CommandInterface) ManagerInterface
	GetCommand(string) CommandInterface
	GetCommands() map[string]CommandInterface
	GetName() string
	GetVersion() string
	Run(...string) error
}
