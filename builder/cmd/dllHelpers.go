package main

import (
	"fmt"
	"io"
)

func writeExportToSB(iw io.Writer, exportName string) error {
	var err error
	if exportName != "init" {
		_, err = fmt.Fprintf(iw, "//export %s\n", exportName)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(iw, "func %s() {\n\t// TODO\n}\n\n", exportName)
	return err
}

func writeDllInit(iw io.Writer) error {
	var err error
	_, err = fmt.Fprint(iw, "func init() {\n\tbin = \"dll\"\n\tmain()\n}\n\n")
	return err
}
