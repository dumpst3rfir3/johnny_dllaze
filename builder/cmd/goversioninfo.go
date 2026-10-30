package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/josephspurrier/goversioninfo"
)

func generateVersioninfo() error {
	var (
		err          error
		tmplBytes    []byte
		buf          bytes.Buffer
		vi           = &goversioninfo.VersionInfo{}
		viConfigFile = "versioninfo.template"
		tmpl         *template.Template
		data         struct {
			Name     string
			FileName string
		}
	)

	tmplBytes, err = os.ReadFile(viConfigFile)
	if err != nil {
		return errors.New("failed to read versioninfo.template: " + err.Error())
	}

	tmpl, err = template.New("versioninfo").Parse(string(tmplBytes))
	if err != nil {
		return errors.New("failed to parse versioninfo template: " + err.Error())
	}

	data.FileName = flags.outputDll
	data.Name = strings.TrimSuffix(flags.outputDll, filepath.Ext(flags.outputDll))

	if err = tmpl.Execute(&buf, data); err != nil {
		return errors.New("failed to render versioninfo template: " + err.Error())
	}

	if err = vi.ParseJSON(buf.Bytes()); err != nil {
		return errors.New("failed to parse versioninfo JSON: " + err.Error())
	}

	vi.Build()

	if err = vi.WriteSyso(
		filepath.Join("dll", "resource_generated.syso"),
		"amd64",
	); err != nil {
		return errors.New("failed to write resource syso: " + err.Error())
	}

	return nil
}
