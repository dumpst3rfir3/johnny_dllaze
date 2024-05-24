# Johnny DLLaze, GOst Writer

![](johnny_dllaze.jpg)

This tool was developed for red team or other security testing purposes. It will simply take a shellcode (e.g., .bin) file, obfuscate the shellcode using [Babble](https://github.com/mjwhitta/babble), and then build a DLL (written in Go) that decodes the shellcode (in place, in memory) and executes it using VirtualAlloc/RtlCopyMemory/VirtualProtect/CreateThread. The DLL is also built with customizable Windows Version Info to make it appear more like a "real" DLL.

Additionally, thanks to [ineffectivecoder](https://github.com/ineffectivecoder), this tool now also generates a "bring your own sideload" executable (goader.exe) that can be used to load the DLL (which will execute the shellcode) by simply placing the DLL in the same directory as goader. Furthermore, the tool can optionally create an ISO file with both goader and the DLL (which will given the hidden attribute within the ISO so that, by default, it won't be viewable in Windows) that can be used for social engineering attack tests.

## Prerequisites

```
go
mingw-w64-gcc
cdrtools
# cdrtools is needed for mkisofs, only if using the option to
# create the iso
```

## Usage

This tool was designed to be used in Linux.

First, edit the versioninfo.json file in the goDLL directory, if desired. You can add your own description, company/copyright info., etc.

Then, from the root directory of the repository, run:

```
./generate_payload.sh /path/to/your/payload.bin [output_dll_filename.dll] [output_iso_filename.dll]
```

If the optional second parameter is passed, that will be used as the output filename for the DLL. Otherwise, it will be named updater.dll. Once it's built, it can be executed with regsvr32, rundll32 (Main, DllRegisterServer, and DllUnRegisterServer are all exported), or sideload (e.g., with goader.exe).

If the optional third parameter is passed, that be used as the output filename for the ISO image. Otherwise, no ISO file will be created.

All output files, including the DLL, the goader.exe file for sideloading, and the optional ISO file, will be located in the `payloads` directory.

## Credit

All the goader sideloader and ISO creation code was written and contributed by [ineffectivecoder](https://github.com/ineffectivecoder). Special shout out to [mjwhitta](https://github.com/mjwhitta) for all of his help, especially the [arTTY](https://github.com/mjwhitta/arTTY) logo used in the generator script.

These great libraries were used in the code:
- [https://github.com/mjwhitta/babble](https://github.com/mjwhitta/babble)
- [https://github.com/josephspurrier/goversioninfo](https://github.com/josephspurrier/goversioninfo)

...and I shamelessly copied code from these sources:
- [https://github.com/mjwhitta/goDLL](https://github.com/mjwhitta/goDLL)
- [https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread](https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread)
