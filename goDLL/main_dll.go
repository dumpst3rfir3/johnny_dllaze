//go:build dll && windows

package main

import "C"

// DllRegisterServer is an entry point called by regsvr32.exe.
//
//export DllRegisterServer
func DllRegisterServer() {}

// DllUnregisterServer is an entry point called by regsvr32.exe.
//
//export DllUnregisterServer
func DllUnregisterServer() {}

//export Now
func Now() {}

// This causes the real main to be called by LoadLibrary() and
// rundll32.exe, hence why Main() is empty.
func init() {
	bin = "dll"
	main()
}

// Main is an entry point to call when using rundll32.exe.
//
//export Main
func Main() {}
