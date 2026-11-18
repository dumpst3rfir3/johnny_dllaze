# Johnny DLLaze, GOst Writer

![](johnny_dllaze.jpg)

This tool was developed for red team or other security testing purposes. It
will simply take a shellcode (e.g., .bin) file, obfuscate the shellcode using
[Babble](https://github.com/mjwhitta/babble), and then build a DLL (written in
Go) that decodes the shellcode (in place, in memory) and executes it using
VirtualAlloc/RtlCopyMemory/VirtualProtect/CreateThread. The DLL is also built
with customizable Windows
[versioninfo](https://learn.microsoft.com/en-us/windows/win32/menurc/versioninfo-resource)
to make it appear more like a "real" DLL. Optionally, the tool supports
adaptive DLL hijacking (inspired by
[Koppeling](https://github.com/TheWover/Koppeling)).

Additionally, thanks to
[ineffectivecoder](https://github.com/ineffectivecoder), this tool also
generates a "bring your own sideload" executable (goader.exe) that can be used
to load the DLL (which will execute the shellcode) by simply placing the DLL
in the same directory as goader. Furthermore, the tool can optionally create
an ISO file with both goader and the DLL that can be used for social
engineering attack tests.

## Prerequisites

- [Go](https://go.dev/) 1.26+
- [MinGW-w64](https://www.mingw-w64.org/downloads) (`gcc` on Windows,
  `x86_64-w64-mingw32-gcc` on Linux/macOS) - required for CGO DLL compilation
- `goversioninfo` - auto-installed by the builder if not present

## Usage

The builder works on Linux, macOS, and Windows.

### Customizing Versioninfo

Before building, you can edit the versioninfo files to control how the output
binaries appear on Windows:

- `versioninfo.template` (project root) - controls the DLL's Windows version
  info. Supports `{{.Name}}` and `{{.FileName}}` template variables which are
  filled in from the output DLL name.
- `goader/versioninfo.json` - controls the sideloader EXE's Windows version
  info.

### Running

From the root directory of the repository, you can run directly without a
separate build step:

```
go run ./builder/cmd/ -s /path/to/payload.bin [OPTIONS]
```

Or build a standalone binary first (optional):

```
go build -o johnny_dllaze ./builder/cmd/
```

Then run it:

```
# Linux/macOS
./johnny_dllaze -s /path/to/payload.bin [OPTIONS]

# Windows
johnny_dllaze.exe -s /path/to/payload.bin [OPTIONS]
```

Key options:

| Flag | Long form | Default | Description |
|------|-----------|---------|-------------|
| `-s` | `--shellcodeFile` | *(required)* | Path to shellcode `.bin` file |
| `-o` | `--output` | `updater.dll` | Output DLL filename |
| `-E` | `--outputExe` | `goader.exe` | Output EXE filename |
| `-I` | `--outputIso` | *(none)* | Output ISO filename (omit to skip ISO creation) |
| `-r` | `--regsvr32` | false | Export `DllRegisterServer`/`DllUnregisterServer` for use with `regsvr32` |
| `-d` | `--dll` | *(none)* | Path to a legitimate DLL to mimic; enables adaptive hijacking mode where all original exports are forwarded |
| `-f` | `--function` | *(none)* | Name of an extra exported function to add. In standard mode, enables `rundll32` execution with the given name. In adaptive hijacking mode, stubs this export instead of forwarding it to the original DLL. If omitted in either mode, no extra export is added |

Run `go run ./builder/cmd/ --help` for the full list.

All output files (DLL, goader.exe, and optional ISO) are placed in the
`payloads/` directory.

## Important Notes

- **Default execution method (standard mode)**: By default, the DLL has no
  exported functions. The shellcode executes from `DllMain` when the DLL is
  loaded, so it can only be run via goader (or any other tool that calls
  `LoadLibrary`). To also support `rundll32`, pass `-f <FuncName>`; to support
  `regsvr32`, pass `-r`.
- **ISO**: The ISO is generated using a pure Go library. The DLL is given the
  hidden attribute within the ISO, matching the behavior of the previous
  version of this tool.
- **Adaptive hijacking - how execution works**: The shellcode executes from
  `DllMain` (via Go's `init()`) when the DLL is loaded, not from a specific
  exported function. All original exports are forwarded to the legitimate DLL
  so the host application keeps running normally. The `-f` flag does not
  change when the shellcode runs; it only controls which export becomes a
  no-op stub instead of forwarding to the original DLL (useful if a particular
  export conflicts with your payload).
- **Adaptive hijacking - loader lock caveat**: Because the shellcode thread is
  started from `DllMain`, calling `WaitForSingleObject` with an infinite
  timeout (as this tool does) can deadlock if the shellcode itself calls
  `LoadLibrary`, since that call will try to acquire the loader lock that
  `DllMain` is already holding. This is typically fine for shellcode that
  doesn't load additional DLLs, but may hang for reflective loaders or other
  payloads that do.

## Credit

All the goader sideloader and ISO creation code was written and contributed by
[ineffectivecoder](https://github.com/ineffectivecoder). Special shout out to
[mjwhitta](https://github.com/mjwhitta) for all of his help, especially the
[arTTY](https://github.com/mjwhitta/arTTY) logo used in the generator script.

These great libraries were used in the code:
- https://github.com/mjwhitta/babble
- https://github.com/josephspurrier/goversioninfo
- https://github.com/kdomanski/iso9660

The adaptive DLL hijacking technique is inspired by
[TheWover/Koppeling](https://github.com/TheWover/Koppeling). See the original
[NetSPI blog post](https://www.netspi.com/adaptive-dll-hijacking/) for
background.

...and I shamelessly copied code from these sources:
- https://github.com/mjwhitta/goDLL
- https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread
