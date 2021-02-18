// author: wsfuyibing <websearch@163.com>
// date: 2021-02-17

// 命令: 创建项目配置文件.
package kv

import (
	"regexp"
	"time"

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
	// 创建实例.
	o := new(command)
	o.Initialize()
	o.SetName("kv")
	o.SetDescription("创建项目配置文件")
	// 加入选项
	o.AddOption(
		base.NewOption("addr", base.OptionModeRequired, base.OptionValueModeString).SetShortName("a").SetDescription("Consul地址"),
		base.NewOption("name", base.OptionModeRequired, base.OptionValueModeString).SetShortName("n").SetDescription("Consul中注册的Key名称"),
		base.NewOption("path", base.OptionModeOptional, base.OptionValueModeString).SetShortName("p").SetDefaultValue("./config").SetDescription("配置文件目录, 默认: ./config"),
		base.NewOption("scheme", base.OptionModeOptional, base.OptionValueModeString).SetShortName("s").SetDefaultValue("http").SetDescription("协议名称, 可选: http, https, 默认: http"),
		base.NewOption("timeout", base.OptionModeOptional, base.OptionValueModeInteger).SetDefaultValue(2).SetDescription("超时时长, 单位: 秒"),
		base.NewOption("upload", base.OptionModeOptional, base.OptionValueModeNone).SetShortName("u").SetDescription("上传初始配置至Consul"),
	)
	return o
}

// Parse arguments.
func (o *command) Run(manager base.ManagerInterface, args []string) error {
	var err error
	var ok bool
	var opt base.OptionInterface
	var timeout = 0
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
	// config: use timeout.
	if opt, err = o.GetOption("timeout"); err != nil {
		return err
	}
	if timeout, err = opt.ToInt(); err != nil {
		return err
	}
	cfg.WaitTime = time.Duration(timeout) * time.Second
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
		return err
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
		files: make(map[string][]string),
	}).run(name)
}
