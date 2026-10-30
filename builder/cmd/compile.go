package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/kdomanski/iso9660"
)

func compileAll() error {
	var err error

	if err = os.MkdirAll("payloads", 0o755); err != nil {
		return errors.New("failed to create payloads directory: " + err.Error())
	}

	if err = compileDLL(); err != nil {
		return err
	}

	if err = compileEXE(); err != nil {
		return err
	}

	if flags.outputIsoFile != "" {
		if err = generateISO(); err != nil {
			return err
		}
	}

	return nil
}

func cc() string {
	if runtime.GOOS == "windows" {
		return "gcc"
	}
	return "x86_64-w64-mingw32-gcc"
}

func compileDLL() error {
	var (
		err error
		cmd *exec.Cmd
		out []byte
	)

	cmd = exec.Command(
		"go", "build",
		"--buildmode=c-shared",
		"--buildvcs=false",
		"--ldflags=-s -w",
		"-o", filepath.Join("..", "payloads", flags.outputDll),
		"--tags=dll",
		"--trimpath",
		".",
	)
	cmd.Dir = "dll"
	cmd.Env = append(os.Environ(),
		"CC="+cc(),
		"CGO_ENABLED=1",
		"GOOS=windows",
		"GOARCH=amd64",
	)

	if out, err = cmd.CombinedOutput(); err != nil {
		return errors.New(
			"failed to build DLL: " + err.Error() + "\n" + string(out),
		)
	}

	log.Printf("[+] Built payloads/%s", flags.outputDll)
	return nil
}

func compileEXE() error {
	var (
		err error
		cmd *exec.Cmd
		out []byte
	)

	cmd = exec.Command("goversioninfo", "-64")
	cmd.Dir = "goader"
	if out, err = cmd.CombinedOutput(); err != nil {
		return errors.New(
			"failed to generate EXE versioninfo: " + err.Error() + "\n" + string(out),
		)
	}

	cmd = exec.Command(
		"go", "build",
		"--buildvcs=false",
		"--ldflags=-s -w -X main.dll="+flags.outputDll,
		"-o", filepath.Join("..", "payloads", flags.outputExe),
		"--trimpath",
		".",
	)
	cmd.Dir = "goader"
	cmd.Env = append(os.Environ(),
		"GOOS=windows",
		"GOARCH=amd64",
	)

	if out, err = cmd.CombinedOutput(); err != nil {
		return errors.New(
			"failed to build EXE: " + err.Error() + "\n" + string(out),
		)
	}

	log.Printf("[+] Built payloads/%s", flags.outputExe)
	return nil
}

func generateISO() error {
	var (
		err     error
		writer  *iso9660.ImageWriter
		entries []fs.DirEntry
		f       *os.File
		output  *os.File
	)

	writer, err = iso9660.NewWriter()
	if err != nil {
		return errors.New("failed to create ISO writer: " + err.Error())
	}
	defer writer.Cleanup()

	entries, err = os.ReadDir("payloads")
	if err != nil {
		return errors.New("failed to read payloads directory: " + err.Error())
	}

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == flags.outputIsoFile {
			continue
		}
		f, err = os.Open(filepath.Join("payloads", entry.Name()))
		if err != nil {
			return errors.New("failed to open " + entry.Name() + ": " + err.Error())
		}
		err = writer.AddFile(f, entry.Name())
		f.Close()
		if err != nil {
			return errors.New("failed to add " + entry.Name() + " to ISO: " + err.Error())
		}
	}

	output, err = os.OpenFile(
		filepath.Join("payloads", flags.outputIsoFile),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0o644,
	)
	if err != nil {
		return errors.New("failed to create ISO file: " + err.Error())
	}
	defer output.Close()

	if err = writer.WriteTo(output, "You've Been GOadered"); err != nil {
		return errors.New("failed to write ISO: " + err.Error())
	}

	log.Printf("[+] ISO generated: payloads/%s", flags.outputIsoFile)
	return nil
}
