package main

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"
)

type Hooks struct {
	OnNewEmail         []*goja.Value
	BeforeSendingEmail []func(e *Email)
}

type MainModule struct {
	Vm    *goja.Runtime
	hooks Hooks
}

type Email struct {
	Subject string
	Body    string
}

func (h Hooks) Init() {
	h.OnNewEmail = make([]*goja.Value, 0)
	h.BeforeSendingEmail = make([]func(e *Email), 0)
}

func (mm MainModule) Init() {
	mm.hooks.Init()
}

func (mm MainModule) RegisterHook(hook string, fn goja.Value) {
	fmt.Println("Register Hook start ")
	//var cb func(e *Email)
	//mm.Vm.ExportTo(fn, &cb)
	mm.hooks.OnNewEmail = append(mm.hooks.OnNewEmail, &fn)
	fmt.Printf("hooks len %d", len(mm.hooks.OnNewEmail))
	//fmt.Println("OnNewEmail call in go ...")
	//cb(&NewEmail{})

}

func (mm *MainModule) TriggerNewEmailEvent(email *Email) {

	fmt.Println("TriggerNewEmailEvent start ")

	fmt.Printf("hooks len %d", len(mm.hooks.OnNewEmail))

	for _, newEmail := range mm.hooks.OnNewEmail {
		fmt.Printf("loop 1")
		var cb func(e *Email)
		mm.Vm.ExportTo(*newEmail, &cb)
		fmt.Printf("Calling call back %v ", newEmail)
		cb(email)
	}
}

func main() {
	vm := goja.New()

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	mod := MainModule{Vm: vm}

	mod.Init()

	vm.Set("myemail", mod)

	script := `
		console.log("JS code started ")
		
		myemail.RegisterHook("onEmailReceived", iGotEmail)
		
		function iGotEmail(newEmail)
		{
			console.log("New Email call back received")
		}
	`
	prg, err := goja.Compile("", script, true)
	if err != nil {
		fmt.Printf("Error compiling the script %v ", err)
		return
	}

	_, err = vm.RunProgram(prg)

	email := &Email{
		Subject: "Your order for blue widgets",
		Body:    " Will be delivered in 1.3 micro seconds",
	}

	mod.TriggerNewEmailEvent(email)
}
