package main

import (
	"github.com/mkrautz/objc/AppKit"
	"github.com/mkrautz/objc/Foundation"
)

func main() {
	rect := foundation.NSRectMake(0, 0, 400, 600)
	window := appkit.NewNSWindow(
		rect,
		appkit.NSTitledWindowMask|appkit.NSClosableWindowMask|appkit.NSMiniaturizableWindowMask,
		appkit.NSBackingStoreBuffered,
		false,
	)
	window.SetTitle("Hello world")
	window.MakeKeyAndOrderFront(window)

	app := appkit.NSSharedApplication()
	app.Run()
}
