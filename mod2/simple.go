package mod2

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

func SimpleJS() {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	script := `
		console.log("Hello world - from Javascript inside Go! ")
	`
	fmt.Println("Compiling ... ")
	prog, err := goja.Compile("", script, true)
	if err != nil {
		fmt.Printf("Error compiling the script %v ", err)
		return
	}
	fmt.Println("Running ... \n ")
	_, err = vm.RunProgram(prog)

}
