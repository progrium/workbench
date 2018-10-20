#!/usr/bin/env osascript -l JavaScript
// http://www.cocoawithlove.com/2010/09/minimalist-cocoa-programming.html
ObjC.import("Cocoa");
ObjC.import("CoreGraphics");

function run(argv) {
    // app = Application.currentApplication();
    // app.includeStandardAdditions = true

    // ObjC.registerSubclass({
    //     name: "AppDelegate",
    //     methods: {
    //         "applicationDidFinishLaunching:": {
    //             types: ["void", ["id"]],
    //             implementation: function (notification) {
    //                  console.log("Started!");
    //             }
    //         },
    //         "buttonClicked:": {
    //             types: ["void", ["id"]],
    //             implementation: function (notification) {
    //                  console.log("Clicked!");
    //             }
    //         }
    //     }
    // });

    // var appDelegate = $.AppDelegate.alloc.init;
    var window = makeWindow();
    window.makeKeyAndOrderFront(window);
    var cocoaApp = $.NSApplication.sharedApplication;
    cocoaApp.setActivationPolicy($.NSApplicationActivationPolicyRegular);
    // cocoaApp.delegate = appDelegate;
    // cocoaApp.activateIgnoringOtherApps(true);
    cocoaApp.run;


}

function makeWindow() {
    var rect = $.CGRectMake($(0), $(0), $(400), $(200));
    console.log(rect)
    var window = $.NSWindow.alloc.initWithContentRectStyleMaskBackingDefer(
        rect,
        $.NSTitledWindowMask | $.NSClosableWindowMask | $.NSMiniaturizableWindowMask,
        $.NSBackingStoreBuffered,
        false
    );
                    
    window.center;
    window.title = "Choose and Display Image";
    return window;
}
