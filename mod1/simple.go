package mod1

import (
	"fmt"

	"github.com/dop251/goja"
)

func SimpleJS() {
	vm := goja.New()

	script := `
		var message = "a simple global message"
	`

	fmt.Println("Compiling ... ")
	prog, err := goja.Compile("", script, true)
	if err != nil {
		fmt.Printf("Error compiling the script %v ", err)
		return
	}
	fmt.Println("Running ... \n ")
	_, err = vm.RunProgram(prog)
	if err != nil {
		fmt.Printf("Error running the script %v ", err)
		return
	}
	msg := vm.Get("message")
	fmt.Printf("Result from script: %s\n", msg.ToString())
}
