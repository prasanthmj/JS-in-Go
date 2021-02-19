# Just to illustrate Running custom Javascript in Go
Uses GoJa library to compile and run Javascript in pure native Go.

* mod1 - set a global variable in the script and get it in Go
* mod2 - make the script write a message to the console 
* mod3 - Call a function in the script from Go
* mod4 - Make a global object available in the script. Call functions on the object from the script. Trigger customizable events from the Go app.
* mod5 - Trigger custom events from Go, pass custom object to the callback function in Javascript, make it possible to call functions on the custom object