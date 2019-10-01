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
	travisTmpl.Execute(os.Stdout, travisData{
		Decrypt:   decryptCommand,
		Binary:    binary,
		Version:   version,
		Project:   project,
		SecretKey: filepath.Base(secretKey),
	})
	return nil
}

type travisData struct {
	Decrypt   string
	Binary    string
	Version   string
	Project   string
	SecretKey string
}

var travisTmpl, _ = template.New("travis").Parse(`branches:
  only:
    - master

before_install:
 - {{.Decrypt}}

install:
 - wget https://storage.googleapis.com/jbd-releases/{{.Binary}}-{{.Version}} && chmod +x ./{{.Binary}}-{{.Version}}

# Build the website
script:
 - GOOGLE_APPLICATION_CREDENTIALS={{.SecretKey}} ./{{.Binary}}-{{.Version}} -project={{.Project}}
`)
