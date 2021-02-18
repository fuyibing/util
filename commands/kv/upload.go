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

	"github.com/fuyibing/util/commands/base"
)

// Upload KV to consul.
type uploadKv struct {
	cmd     base.CommandInterface
	path    string
	cli     *api.Client
	content string
}

func (o *uploadKv) run(key string) error {
	var err error
	var fs []os.FileInfo
	var fp *os.File
	var pair *api.KVPair
	// read fail.
	fi := 0
	fs, err = ioutil.ReadDir(o.path)
	if err != nil {
		return err
	}
	// exist check.
	pair, _, err = o.cli.KV().Get(key, nil)
	if err != nil {
		return err
	}
	if pair != nil {
		return errors.New("key exists")
	}
	// open directory.
	for _, f := range fs {
		// continue if directory.
		if f.IsDir() {
			continue
		}
		// continue if not yaml.
		if !regexpIsYamlFile.MatchString(f.Name()) {
			continue
		}
		// head.
		o.content += fmt.Sprintf("# [%s]\n", f.Name())
		o.content += fmt.Sprintf("%s: \n", f.Name())
		// open file.
		if fp, err = os.Open(o.path + "/" + f.Name()); err != nil {
			return err
		}
		// loop line.
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
		_ = fp.Close()
		o.content += "\n"
	}
	// no file found.
	if fi == 0 {
		return errors.New(fmt.Sprintf("Command %s: no yaml config file found: %s", o.cmd.GetName(), o.path))
	}
	// key pair.
	pair = &api.KVPair{Key: key}
	pair.ModifyIndex = 0
	pair.Value = []byte(o.content)
	// put key.
	_, err = o.cli.KV().Put(pair, nil)
	if err != nil {
		return err
	}
	return nil
}
