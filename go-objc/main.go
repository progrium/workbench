package main

import (
	"fmt"

	"github.com/progrium/objc/AppKit"
	foundation "github.com/progrium/objc/Foundation"
)

func main() {
	window := appkit.NewNSWindow(
		foundation.NSRect{
			10.0, 10.0, 500.0, 400.0,
		},
		appkit.NSTitledWindowMask|appkit.NSClosableWindowMask|appkit.NSMiniaturizableWindowMask,
		appkit.NSBackingStoreBuffered,
		false,
	)
	window.SetTitle("Hello world")
	window.MakeKeyAndOrderFront(window)

	app := appkit.NSSharedApplication()
	app.SendMsg("setActivationPolicy:", 0)
	fmt.Println("running...")
	app.Run()
}

//statusBarItem := objc.GetClass("NSStatusBar").SendMsg("systemStatusBar").SendMsg("statusItemWithLength:", -1.0)

//statusBarItem.SendMsg("button").SendMsg("setTitle:", foundation.NSStringFromString("Hello world").Object)
