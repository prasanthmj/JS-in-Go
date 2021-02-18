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

type Email struct {
	Subject string
	Body    string
}

func (h Hooks) Init() {
	h.OnNewEmail = make([]*goja.Value, 0)
	h.BeforeSendingEmail = make([]func(e *Email), 0)
}

func (h *Hooks) TriggerNewEmailEvent(email *Email, vm *goja.Runtime) {

	fmt.Println("TriggerNewEmailEvent start ")

	fmt.Printf("hooks len %d\n", len(h.OnNewEmail))

	for _, newEmail := range h.OnNewEmail {
		var cb func(e *Email)
		vm.ExportTo(*newEmail, &cb)
		cb(email)
	}
}

func main() {
	var hooks Hooks
	hooks.Init()

	vm := goja.New()

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	obj := vm.NewObject()

	obj.Set("RegisterHook", func(hook string, fn goja.Value) {
		fmt.Printf("RegisterHook called ")
		hooks.OnNewEmail = append(hooks.OnNewEmail, &fn)
	})

	vm.Set("myemail", obj)

	script := `
		console.log("JS code started ")
		
		myemail.RegisterHook("onEmailReceived", iGotEmail)
		
		function iGotEmail(newEmail)
		{
			console.log("New Email call back received. email=",newEmail.Subject, newEmail.Body)
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

	fmt.Println(" Triggering event ")
	hooks.TriggerNewEmailEvent(email, vm)
}
