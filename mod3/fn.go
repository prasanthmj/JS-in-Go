package mod3

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

func CallJSFunction() {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	script := `
		function myFunction(param)
		{
			console.log("myFunction running ...")
			console.log("Param = ", param)
			return "Nice meeting you, Go"
		}
	`

	prog, err := goja.Compile("", script, true)
	if err != nil {
		fmt.Printf("Error compiling the script %v ", err)
		return
	}
	_, err = vm.RunProgram(prog)

	var myJSFunc goja.Callable
	err = vm.ExportTo(vm.Get("myFunction"), &myJSFunc)
	if err != nil {
		fmt.Printf("Error exporting the function %v", err)
		return
	}

	res, err := myJSFunc(goja.Undefined(), vm.ToValue("message from go"))
	if err != nil {
		fmt.Printf("Error calling function %v", err)
		return
	}
	fmt.Printf("Returned value from JS function\n%s \n", res.ToString())

}
