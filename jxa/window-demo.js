#!/usr/bin/env osascript -l JavaScript
// http://www.cocoawithlove.com/2010/09/minimalist-cocoa-programming.html
ObjC.import("Cocoa");

function run(argv) {
    app = Application.currentApplication();
    app.includeStandardAdditions = true

    ObjC.registerSubclass({
        name: "AppDelegate",
        methods: {
            "applicationDidFinishLaunching:": {
                types: ["void", ["id"]],
                implementation: function (notification) {
                     console.log("Started!");
                }
            },
            "buttonClicked:": {
                types: ["void", ["id"]],
                implementation: function (notification) {
                     console.log("Clicked!");
                }
            }
        }
    });

    var appDelegate = $.AppDelegate.alloc.init;
    var window = makeWindow(appDelegate);
    window.makeKeyAndOrderFront(window);
    var cocoaApp = $.NSApplication.sharedApplication;
    cocoaApp.setActivationPolicy($.NSApplicationActivationPolicyRegular);
    cocoaApp.delegate = appDelegate;
    cocoaApp.activateIgnoringOtherApps(true);
    cocoaApp.run;


}

function makeWindow(appDelegate) {
    var styleMask = $.NSTitledWindowMask | $.NSClosableWindowMask | $.NSMiniaturizableWindowMask;
    var windowHeight = 400;
    var windowWidth = 600;
    var ctrlsHeight = 80;
    var minWidth = 400;
    var minHeight = 340;
    var window = $.NSWindow.alloc.initWithContentRectStyleMaskBackingDefer(
        $.NSMakeRect(0, 0, windowWidth, windowHeight),
        styleMask,
        $.NSBackingStoreBuffered,
        false
    );
                    
    var btn = $.NSButton.alloc.initWithFrame($.NSMakeRect(25, (windowHeight - 100), 200, 25));
    btn.title = "Update Label";
    btn.bezelStyle = $.NSRoundedBezelStyle;
    btn.buttonType = $.NSMomentaryLightButton;
    btn.target = appDelegate;
    btn.action = "buttonClicked:";
    
    window.contentView.addSubview(btn);
    
    window.center;
    window.title = "Choose and Display Image";
    return window;
}
