package main
//go:generate goversioninfo -64
import (
	"fmt"

	"golang.org/x/sys/windows"
)

var dll string

func main() {
	if _, e := windows.LoadLibrary(dll); e != nil {
		println(e.Error())
	}

	fmt.Println("Press enter to exit")
	fmt.Scanln()
}
