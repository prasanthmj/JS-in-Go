package mod4

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

	for _, newEmail := range h.OnNewEmail {
		var newEmailCallBack func(e *Email)
		vm.ExportTo(*newEmail, &newEmailCallBack)
		newEmailCallBack(email)
	}
}

func WorkWithHooks() {
	var hooks Hooks
	hooks.Init()

	vm := goja.New()

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	obj := vm.NewObject()

	obj.Set("RegisterHook", func(hook string, fn goja.Value) {
		hooks.OnNewEmail = append(hooks.OnNewEmail, &fn)
		fmt.Println("Registered the Hook ")
	})

	vm.Set("myemail", obj)

	script := `
		console.log("JS code started ")
		
		myemail.RegisterHook("onEmailReceived", iGotEmail)
		
		function iGotEmail(newEmail)
		{
			console.log("New Email callback received. \n",newEmail.Subject,"\n", newEmail.Body)
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
		Body:    "Will be delivered in 1.3 micro seconds",
	}

	fmt.Println("Triggering the event ")
	hooks.TriggerNewEmailEvent(email, vm)
}
