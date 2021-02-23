// author: wsfuyibing <websearch@163.com>
// date: 2021-02-18

package kv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/consul/api"

	"github.com/fuyibing/util/cmds/base"
)

// Upload KV to consul.
type uploadKv struct {
	cmd     base.CommandInterface
	path    string
	cli     *api.Client
	content string
}

// Run upload.
func (o *uploadKv) run(key string) error {
	var err error
	var fi = 0
	var fs []os.FileInfo
	var fp *os.File
	var pair *api.KVPair
	// 1. read directory error.
	if fs, err = ioutil.ReadDir(o.path); err != nil {
		return errors.New(fmt.Sprintf("Command %s: open directory error: %v", o.cmd.GetName(), o.path))
	}
	// 2. open remote key error.
	if pair, _, err = o.cli.KV().Get(key, nil); err != nil {
		return errors.New(fmt.Sprintf("Command %s: open remote key error: %v", o.cmd.GetName(), err))
	}
	if pair != nil {
		return errors.New(fmt.Sprintf("Command %s: remote key exist already: %v", o.cmd.GetName(), key))
	}
	// 3. seek directory.
	for _, f := range fs {
		// 3.1 ignore if is directory.
		if f.IsDir() {
			continue
		}
		// 3.2 ignore if not YAML file.
		if !regexpIsYamlFile.MatchString(f.Name()) {
			continue
		}
		// 3.3 file meta data.
		o.content += fmt.Sprintf("# [%s]\n", f.Name())
		o.content += fmt.Sprintf("%s: \n", f.Name())
		// 3.4 return if open file error.
		if fp, err = os.Open(o.path + "/" + f.Name()); err != nil {
			return err
		}
		o.cmd.Info("Command %s: upload %s", o.cmd.GetName(), f.Name())
		// 3.5 read line by line.
		fi++
		br := bufio.NewReader(fp)
		for {
			b, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			s := string(b)
			// ignore comment.
			if regexpIsComment.MatchString(s) {
				continue
			}
			// ignore empty line.
			if regexpIsEmptyLine.MatchString(s) {
				continue
			}
			// append line.
			o.content += fmt.Sprintf("  %s\n", s)
		}
		// 3.6 close file.
		_ = fp.Close()
		o.content += "\n"
	}
	// no file found.
	if fi == 0 {
		return errors.New(fmt.Sprintf("Command %s: no yaml config file found: %s", o.cmd.GetName(), o.path))
	}
	// key pair.
	pair = &api.KVPair{Key: key}
	pair.Value = []byte(o.content)
	// put key.
	_, err = o.cli.KV().Put(pair, nil)
	if err != nil {
		return err
	}
	return nil
}
