package mod5

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
	Subject  string
	Body     string
	Priority int
	To       []string
}

func (h Hooks) Init() {
	h.OnNewEmail = make([]*goja.Value, 0)
	h.BeforeSendingEmail = make([]func(e *Email), 0)
}

func (h *Hooks) TriggerNewEmailEvent(email *Email, vm *goja.Runtime) {

	eobj := makeEmailJSObject(vm, email)

	for _, newEmail := range h.OnNewEmail {
		var newEmailCallBack func(*goja.Object)
		vm.ExportTo(*newEmail, &newEmailCallBack)
		newEmailCallBack(eobj)
	}
}

func makeEmailJSObject(vm *goja.Runtime, email *Email) *goja.Object {
	obj := vm.NewObject()
	obj.Set("subject", email.Subject)
	obj.Set("body", email.Body)
	obj.Set("to", email.To)
	obj.Set("reply", func(body string) {
		fmt.Printf("Replying:\n%s\n", body)
	})
	obj.Set("setPriority", func(p int) {
		fmt.Printf("Set email priority to %d\n", p)
	})

	obj.Set("moveTo", func(folder string) {
		fmt.Printf("Moving email to folder %s\n", folder)
	})

	return obj
}

func WorkWithHooks() {
	var hooks Hooks
	hooks.Init()

	vm := goja.New()

	new(require.Registry).Enable(vm)
	console.Enable(vm)

	obj := vm.NewObject()

	obj.Set("RegisterHook", func(hook string, fn goja.Value) {
		switch hook {
		case "onEmailReceived":
			hooks.OnNewEmail = append(hooks.OnNewEmail, &fn)
			fmt.Println("Registered onEmailReceived Hook ")
		}

	})

	vm.Set("myemail", obj)

	script := `
		console.log("JS code started ")
		
		myemail.RegisterHook("onEmailReceived", iGotEmail)

		function iGotEmail(newEmail)
		{
			console.log("newEmail, subject %s ", newEmail.subject,  newEmail.to)
			
			if(newEmail.subject.startsWith("URGENT:"))
			{
				newEmail.setPriority(5)
				
				newEmail.reply("Hello,\n Received your email. We will respond on priority basis. \n\nThanks\n")
				return
			}
			else if(newEmail.to.includes("sales@website"))
			{
				newEmail.moveTo("Sales")
				return
			}
		}
	`
	prg, err := goja.Compile("", script, true)
	if err != nil {
		fmt.Printf("Error compiling the script %v ", err)
		return
	}

	_, err = vm.RunProgram(prg)

	email1 := &Email{
		Subject: "URGENT: Systems down!",
		Body:    "5 of your systems are down at the moment",
		To:      []string{"some@one.cc"},
	}

	fmt.Println("Triggering the urgent email event ")
	hooks.TriggerNewEmailEvent(email1, vm)

	email2 := &Email{
		Subject: "New order received!",
		Body:    "You got new order for 1k blue widgets!",
		To:      []string{"sales@website"},
	}

	fmt.Println("Triggering the sales email event ")
	hooks.TriggerNewEmailEvent(email2, vm)
}
