package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/josephspurrier/goversioninfo"
)

func generateVersioninfo() error {
	var (
		buf          bytes.Buffer
		data         struct {
			Name     string
			FileName string
		}
		err          error
		tmpl         *template.Template
		tmplBytes    []byte
		vi           = &goversioninfo.VersionInfo{}
		viConfigFile = "versioninfo.template"
	)

	tmplBytes, err = os.ReadFile(viConfigFile)
	if err != nil {
		return errors.New("failed to read versioninfo.template: " + err.Error())
	}

	tmpl, err = template.New("versioninfo").Parse(string(tmplBytes))
	if err != nil {
		return fmt.Errorf("failed to parse versioninfo template: %w", err)
	}

	data.FileName = flags.outputDll
	data.Name = strings.TrimSuffix(
		flags.outputDll, filepath.Ext(flags.outputDll),
	)

	if err = tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to render versioninfo template: %w", err)
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

	if err = patchSysoMachineAMD64(
		filepath.Join("dll", "resource_generated.syso"),
	); err != nil {
		return errors.New("failed to patch resource syso: " + err.Error())
	}

	return nil
}

// patchSysoMachineAMD64 fixes the Machine field in a goversioninfo-generated
// COFF resource file from IMAGE_FILE_MACHINE_UNKNOWN (0x0000) to
// IMAGE_FILE_MACHINE_AMD64 (0x8664). goversioninfo writes Machine=0 which
// Go's internal linker accepts but GNU ld (MinGW) rejects.
func patchSysoMachineAMD64(path string) error {
	var (
		b   []byte
		err error
	)

	b, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(b) < 2 {
		return errors.New("syso file too small")
	}

	b[0] = 0x64
	b[1] = 0x86

	return os.WriteFile(path, b, 0o600)
}
