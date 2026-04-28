package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mjwhitta/cli"
)

var flags struct {
	adaptHijackDll      string
	adaptHijackFunction string
	keyFile             string
	outputDll           string
	outputExe           string
	outputIsoFile       string
	regSvr32            bool
	runDll32Func        string
	shellcodeFile       string
}

func init() {
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
		"hijack, if adaptive DLL hijacking is desired (like Koppeling).")
	cli.Flag(&flags.adaptHijackFunction, "f", "function", "", "Function to "+
		"hijack in the DLL, if using adaptive DLL hijacking.")
	cli.Flag(&flags.keyFile, "k", "keyFile", "", "Path to the key file "+
		"that Babble will use for shellcode obfuscation. If not provided, "+
		"one will be generated for you.")
	cli.Flag(&flags.outputDll, "D", "outputDll", "updater.dll", "Name of the "+
		"output DLL (default: updater.dll)")
	cli.Flag(&flags.outputExe, "E", "outputExe", "goader.exe", "Name of the "+
		"output executable (default: goader.exe)")
	cli.Flag(&flags.outputIsoFile, "I", "outputIsoFile", "", "Name of the "+
		"output ISO file to use, if using ISO-based delivery (e.g., for "+
		"HTML smuggling - by default, none will be generated)")
	cli.Flag(&flags.regSvr32, "r", "regSvr32", false, "If set, the DLL will "+
		"be generated with the DllRegisterServer and DllUnregisterServer "+
		"entry points so that it can be executed with regsvr32")
	cli.Flag(&flags.runDll32Func, "R", "runDll32Func", "", "If set, the DLL "+
		"the DLL will be generated with an exported function with the "+
		" specified name, which can be executed with rundll32. If using "+
		"adaptive DLL hijacking, this will be ignored and the hijacked"+
		" function will be used instead. If adaptive DLL hijacking is not "+
		"being used, and no runDll32Func is specified, then a default "+
		"function named \"Dllaze\" will be used")
	cli.Flag(&flags.shellcodeFile, "s", "shellcodeFile", "", "Path to the "+
		"shellcode file (REQUIRED)")

	cli.Parse()
	if flags.shellcodeFile == "" {
		fmt.Println("[!] Error: shellcode file is required")
		cli.Usage(1)
	}
	if flags.adaptHijackDll == "" &&
		flags.runDll32Func == "" &&
		!flags.regSvr32 {
		flags.runDll32Func = "Dllaze"
	}
	if flags.adaptHijackDll == "" && flags.adaptHijackFunction != "" {
		fmt.Println("[!] Error: if adaptHijackFunction is specified, " +
			"adaptHijackDll must also be specified")
		cli.Usage(1)
	}
}
