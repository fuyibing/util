// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 创建项目配置文件.
package kv

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/consul/api"

	"github.com/fuyibing/util/commands/base"
)

var (
	regexpIsComment   = regexp.MustCompile(`^\s*#`)
	regexpIsEmptyLine = regexp.MustCompile(`^\s+$`)
	regexpDepth       = regexp.MustCompile(`kv://([/_a-zA-Z0-9]+)`)
	regexpIsYamlFile  = regexp.MustCompile(`^([_a-zA-Z0-9\-]+).yaml\s*[:]*\s*$`)
)

type command struct {
	base.Command
}

func New() base.CommandInterface {
	o := new(command)
	o.Initialize()
	o.SetName("kv")
	o.SetDescription("Create application config files use consul")
	o.AddOption(
		base.NewOption("addr", base.OptionModeRequired, base.OptionValueModeString).SetShortName("a").SetDescription("Consul address, eg: 192.168.1.1:8500"),
		base.NewOption("name", base.OptionModeRequired, base.OptionValueModeString).SetShortName("n").SetDescription("Consul key name, eg: app/config"),
		base.NewOption("path", base.OptionModeOptional, base.OptionValueModeString).SetShortName("p").SetDefaultValue("./tmp").SetDescription("Config file directory name (default: ./tmp)"),
		base.NewOption("scheme", base.OptionModeOptional, base.OptionValueModeString).SetShortName("s").SetDefaultValue("http").SetDescription("Consul scheme (accept: http|https, default: http)"),
		base.NewOption("upload", base.OptionModeOptional, base.OptionValueModeNone).SetShortName("u").SetDescription("Upload local config files data to remote"),
	)
	return o
}

// Parse arguments.
func (o *command) Run(manager base.ManagerInterface, args []string) error {
	var err error
	var ok bool
	var opt base.OptionInterface
	// argument parse.
	if err = o.ParseArguments(args); err != nil {
		return err
	}
	// get specified option.
	if opt, err = o.GetOption("upload"); err != nil {
		return err
	}
	// to boolean error.
	if ok, err = opt.ToBool(); err != nil {
		return err
	}
	// config: init.
	cfg := api.DefaultConfig()
	// config: use scheme.
	if opt, err = o.GetOption("scheme"); err != nil {
		return err
	}
	if cfg.Scheme, err = opt.ToString(); err != nil {
		return err
	}
	// config: use addr.
	if opt, err = o.GetOption("addr"); err != nil {
		return err
	}
	if cfg.Address, err = opt.ToString(); err != nil {
		return err
	}
	// build: client
	var cli *api.Client
	if cli, err = api.NewClient(cfg); err != nil {
		return errors.New(fmt.Sprintf("Command %s: create consul client error: %v", o.GetName(), err))
	}
	// build: name & path.
	var name, path = "", ""
	if opt, err = o.GetOption("name"); err != nil {
		return err
	}
	if name, err = opt.ToString(); err != nil {
		return err
	}
	if opt, err = o.GetOption("path"); err != nil {
		return err
	}
	if path, err = opt.ToString(); err != nil {
		return err
	}
	if err = os.Mkdir(path, os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("Command %s: create config path error: %v", o.GetName(), err))
		}
	}
	// upload.
	if ok {
		return (&uploadKv{
			cmd:  o,
			cli:  cli,
			path: path,
		}).run(name)
	}
	// download.
	return (&downloadKv{
		cmd:   o,
		cli:   cli,
		path:  path,
	}).run(name)
}
