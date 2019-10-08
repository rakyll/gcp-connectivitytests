// Copyright 2019 Google Inc. All Rights Reserved.
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

package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func generate(t string) error {
	switch t {
	case "travis":
		return generateTravis()
	case "circleci":
		return generateCircle()
	}
	return fmt.Errorf("unsupported gen target %q", t)
}

func generateTravis() error {
	if secretKey == "" {
		return errors.New("-secretkey cannot be empty; provide a Google Cloud secret key")
	}
	if _, err := exec.LookPath("travis"); err != nil {
		return errors.New("travis command is not installed, see https://github.com/travis-ci/travis.rb")
	}
	cmd := exec.Command("travis",
		"encrypt-file",
		secretKey,
		"-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}

	var decryptCommand string
	buf := bytes.NewBuffer(out)
	for {
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if bytes.Contains(line, []byte("openssl aes-256-cbc")) {
			decryptCommand = strings.TrimSpace(string(line))
		}
	}
	return travisTmpl.Execute(os.Stdout, travisData{
		Decrypt:   decryptCommand,
		Binary:    linuxBinary,
		Project:   project,
		SecretKey: filepath.Base(secretKey),
	})
}

func generateCircle() error {
	return circleTmpl.Execute(os.Stdout, travisData{
		Binary:  linuxBinary,
		Project: project,
	})
}

type travisData struct {
	Decrypt   string
	Binary    string
	Project   string
	SecretKey string
}

type circleData struct {
	Binary  string
	Project string
}

var travisTmpl, _ = template.New("travis").Parse(`branches:
  only:
    - master

before_install:
 - {{.Decrypt}}

install:
 - wget https://storage.googleapis.com/jbd-releases/{{.Binary}} && chmod +x ./{{.Binary}}

script:
 - GOOGLE_APPLICATION_CREDENTIALS={{.SecretKey}} ./{{.Binary}} -project={{.Project}}
`)

var circleTmpl, _ = template.New("circleci").Parse(`version: 2
jobs:
  build:
    docker:
      - image: google/cloud-sdk

    steps:
      - run: |
          apt-get install wget -y
          wget https://storage.googleapis.com/jbd-releases/{{.Binary}} && chmod +x ./{{.Binary}}
          echo $GCLOUD_SERVICE_KEY > key.json
          GOOGLE_APPLICATION_CREDENTIALS=key.json ./{{.Binary}} -project={{.Project}}
`)
