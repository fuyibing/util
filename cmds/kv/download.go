// author: wsfuyibing <websearch@163.com>
// date: 2021-02-18

package kv

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/fuyibing/util/cmds/base"
)

// Download KV from consul.
type downloadKv struct {
	cmd      base.CommandInterface
	path     string
	cli      *api.Client
	files    map[string][]string
	origin   bool
	override bool
	content  string
}

// Parse content.
func (o *downloadKv) parseContent() (err error) {
	// catch panic.
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Command %s: panic: %v", o.cmd.GetName(), r))
		}
	}()
	// normal match.
	key := ""
	for _, line := range strings.Split(o.content, "\n") {
		if line == "" || regexpIsEmptyLine.MatchString(line) {
			continue
		}
		if m := regexpIsComment.MatchString(line); m {
			continue
		}
		if m := regexpIsYamlFile.FindStringSubmatch(line); len(m) == 2 {
			key = m[1]
			o.files[key] = make([]string, 0)
			continue
		}
		// append content.
		if key == "" {
			err = errors.New(fmt.Sprintf("Command %s: file name not defined in consul kv", o.cmd.GetName()))
			return
		}
		// key not define.
		if _, ok := o.files[key]; !ok {
			err = errors.New(fmt.Sprintf("Command %s: invalid file name: %s", o.cmd.GetName(), key))
			return
		}
		// append.
		o.files[key] = append(o.files[key], line[2:])
	}
	return
}

// Read value by key name.
func (o *downloadKv) readContent(name string) (string, error) {
	p, _, err := o.cli.KV().Get(name, nil)
	if err != nil {
		return "",
			errors.New(fmt.Sprintf("Command %s: open remote key error: %s", o.cmd.GetName(), err))
	}
	if p == nil {
		return "",
			errors.New(fmt.Sprintf("Command %s: key not found: %s", o.cmd.GetName(), name))
	}
	return string(p.Value), nil
}

// Read depth.
func (o *downloadKv) readDepth() error {
	var err error
	o.content = regexpDepth.ReplaceAllStringFunc(o.content, func(s string) string {
		m := regexpDepth.FindStringSubmatch(s)
		content, e := o.readContent(m[1])
		if e != nil {
			err = e
		}
		return content
	})
	return err
}

func (o *downloadKv) run(key string) error {
	var opt base.OptionInterface
	// download origin.
	opt, _ = o.cmd.GetOption("origin")
	o.origin, _ = opt.ToBool()
	// override if file exist
	opt, _ = o.cmd.GetOption("override")
	o.override, _ = opt.ToBool()
	// loop file.
	for _, name := range strings.Split(key, ",") {
		o.content = ""
		o.files = make(map[string][]string)
		if err := o.runKey(name); err != nil {
			return err
		}
	}
	return nil
}

// Run api.
func (o *downloadKv) runKey(key string) error {
	var err error
	// 1. read base content.
	if o.content, err = o.readContent(key); err != nil {
		return err
	}
	// 2. depth replace.
	if !o.origin {
		if err = o.readDepth(); err != nil {
			return err
		}
	}
	// 3. parse content to file.
	if err = o.parseContent(); err != nil {
		return err
	}
	// 4. write file.
	for name, lines := range o.files {
		if _, err = o.write(key, name, lines); err != nil {
			return err
		}
		o.cmd.Info("Command %s: download %s.", o.cmd.GetName(), name)
	}
	return nil
}

// Write content to file.
func (o *downloadKv) write(key, name string, lines []string) (string, error) {
	var err error
	var f *os.File
	var src = fmt.Sprintf("%s/%s.yaml", o.path, name)
	// check override able if file exists
	if !o.override {
		if f, err = os.OpenFile(src, os.O_RDONLY, os.ModePerm); err == nil {
			_ = f.Close()
			return "", errors.New(fmt.Sprintf("Command %s: file exist: %s", o.cmd.GetName(), name))
		}
	}
	// open and close if end.
	f, err = os.OpenFile(src, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return src, err
	}
	defer func() {
		_ = f.Close()
	}()
	// header
	s := "# config from consul"
	s += fmt.Sprintf("\n# key: %s", key)
	s += fmt.Sprintf("\n# file: %s", src)
	s += fmt.Sprintf("\n# date: %s", time.Now().String())
	// content
	for _, line := range lines {
		s += fmt.Sprintf("\n%s", line)
	}
	// write
	_, err = f.WriteString(s)
	if err != nil {
		return src, err
	}
	return src, nil
}
