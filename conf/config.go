// Copyright 2014 mqant Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gitee.com/yuanxuezhe/dante/log"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	LenStackBuf = 1024

	Conf = Config{}
)

func init() {
	ApplicationDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		ApplicationPath, _ := filepath.Abs(file)
		ApplicationDir, _ = filepath.Split(ApplicationPath)
	}

	defaultConfPath := fmt.Sprintf("%s\\conf\\server.json", ApplicationDir)
	log.Release(defaultConfPath)
	LoadConfig(defaultConfPath)
}

func LoadConfig(Path string) {
	// Read config.
	if err := readFileInto(Path); err != nil {
		panic(err)
	}
}

type Config struct {
	Registermodules  string
	RegisterProtocol string
	RegisterCentor   string
	Log              map[string]interface{}
	Module           map[string][]*ModuleSettings
	Settings         map[string]interface{}
	Mysql            Mysql
}

type Rabbitmq struct {
	Uri          string
	Exchange     string
	ExchangeType string
	Queue        string
	BindingKey   string //
	ConsumerTag  string //消费者TAG
}

type Redis struct {
	Uri   string //redis://:[password]@[ip]:[port]/[db]
	Queue string
}

type Mysql struct {
	Url      string //redis://:[password]@[ip]:[port]/[db]
	Maxcount int
}

type ModuleSettings struct {
	Id string
	//TCPAddr string
	//ProcessID string
	Settings map[string]interface{}
	//Rabbitmq *Rabbitmq
	//Redis    *Redis
}

type SSH struct {
	Host     string
	Port     int
	User     string
	Password string
}

/**
host:port
*/
func (s *SSH) GetSSHHost() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func readFileInto(path string) error {
	var data []byte
	buf := new(bytes.Buffer)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			if len(line) > 0 {
				buf.Write(line)
			}
			break
		}
		if !strings.HasPrefix(strings.TrimLeft(string(line), "\t "), "//") {
			buf.Write(line)
		}
	}
	data = buf.Bytes()
	return json.Unmarshal(data, &Conf)
}

// If read the file has an error,it will throws a panic.
func fileToStruct(path string, ptr *[]byte) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	*ptr = data
}
