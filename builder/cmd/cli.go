package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mjwhitta/cli"
)

var flags struct {
	adaptHijackDll string
	exportFunc     string
	outputDll      string
	outputExe      string
	outputIsoFile  string
	regSvr32       bool
	shellcodeFile  string
}

func init() {
	var err error

	cli.Align = true
	cli.Banner = fmt.Sprintf("%s [OPTIONS]", filepath.Base(os.Args[0]))
	cli.Info(
		"Johnny Dllaze is a red team tool for obfuscating shellcode and " +
			"building a DLL loader to execute it. It will also build an " +
			"executable that will load the DLL (i.e., to \"bring your own " +
			"sideload\"). Optionally, it can also generate an ISO file " +
			"(e.g.,for HTML smuggling delivery), and also supports adaptive " +
			"DLL hijacking (like Koppeling). You can specify which DLL to " +
			"mimic and which function to hijack (i.e., which function will " +
			"actually executed the shellcode).",
	)
	cli.Flag(&flags.adaptHijackDll, "d", "dll", "", "Path to the DLL to "+
		"mimic for adaptive DLL hijacking (like Koppeling).")
	cli.Flag(&flags.exportFunc, "f", "function", "", "Name of an exported "+
		"function to add to the DLL. In standard mode, this enables "+
		"rundll32 execution with the given function name. In adaptive "+
		"hijacking mode (-d), this stubs the named export instead of "+
		"forwarding it to the original DLL. In both modes, if omitted, "+
		"no extra export is added.")
	cli.Flag(&flags.outputDll, "o", "output", "updater.dll", "Name of the "+
		"output DLL (default: updater.dll)")
	cli.Flag(&flags.outputExe, "E", "outputExe", "goader.exe", "Name of the "+
		"output executable (default: goader.exe)")
	cli.Flag(&flags.outputIsoFile, "I", "outputIso", "", "Name of the "+
		"output ISO file to use, if using ISO-based delivery (e.g., for "+
		"HTML smuggling - by default, none will be generated)")
	cli.Flag(&flags.regSvr32, "r", "regsvr32", false, "If set, the DLL will "+
		"be generated with the DllRegisterServer and DllUnregisterServer "+
		"entry points so that it can be executed with regsvr32")
	cli.Flag(&flags.shellcodeFile, "s", "shellcodeFile", "", "Path to the "+
		"shellcode file (REQUIRED)")

	cli.Parse()
	err = validate()
	if err != nil {
		fmt.Printf("[!] Error: %s\n", err)
		cli.Usage(1)
	}

}

func validate() error {
	if flags.shellcodeFile == "" {
		return errors.New("shellcode file is required")
	}

	if flags.adaptHijackDll != "" && flags.outputDll == "updater.dll" {
		flags.outputDll = filepath.Base(flags.adaptHijackDll)
	}

	if strings.Contains(filepath.ToSlash(flags.outputDll), "/") {
		return errors.New("output DLL name cannot contain path separators")
	}
	if strings.Contains(filepath.ToSlash(flags.outputExe), "/") {
		return errors.New("output EXE name cannot contain path separators")
	}
	if flags.outputIsoFile != "" &&
		strings.Contains(filepath.ToSlash(flags.outputIsoFile), "/") {
		return errors.New(
			"output ISO file name cannot contain path separators",
		)
	}

	return nil
}
