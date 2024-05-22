# Johnny DLLaze, GOst Writer

![](johnny_dllaze.jpg)

This tool was developed for red team or other security testing purposes. It will simply take a shellcode (e.g., .bin) file, obfuscate the shellcode using [Babble](https://github.com/mjwhitta/babble), and then build a DLL (written in Go) that decodes the shellcode (in place, in memory) and executes it using VirtualAlloc/RtlCopyMemory/VirtualProtect/CreateThread. The DLL is also built with customizable Windows Version Info to make it appear more like a "real" DLL.

## Prerequisites

```
go
mingw-w64-gcc
```

## Usage

This tool was designed to be used in Linux.

First, edit the versioninfo.json file in the goDLL directory, if desired. You will also probably want to change the name of the output DLL in `goDLL/build_dll_on_linux.sh`. (The default name is updater.dll)

Then, from the root directory of the repository, run:

```
./generate_payload.sh /path/to/your/payload.bin
```

The generated DLL file, updater.dll by default, will be in the goDLL directory.

Once it's built, it can be executed with regsvr32, rundll32 (Main, DllRegisterServer, and DllUnRegisterServer are all exported), or sideload.

## Credit

These great libraries were used in the code:
- [https://github.com/mjwhitta/babble](https://github.com/mjwhitta/babble)
- [https://github.com/josephspurrier/goversioninfo](https://github.com/josephspurrier/goversioninfo)

...and I shamelessly copied code from these sources:
- [https://github.com/mjwhitta/goDLL](https://github.com/mjwhitta/goDLL)
- [https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread](https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread)
