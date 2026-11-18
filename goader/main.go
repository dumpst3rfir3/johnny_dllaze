package main
//go:generate goversioninfo -64
import (
	"fmt"

	"golang.org/x/sys/windows"
)

var dll string

func main() {
	var e error
	_, e = windows.LoadLibrary(dll)
	if e != nil {
		println(e.Error())
	}

	fmt.Println("Press enter to exit")
	fmt.Scanln()
}
