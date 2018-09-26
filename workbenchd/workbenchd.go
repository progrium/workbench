package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/progrium/prototypes/pkg/supervisor"
)

// toDO: rpcdaemon possible halts when reloading a program that exits

func triggerHook(name string, args ...string) {
	path := fmt.Sprintf("hooks/%s", name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}
	cmd := exec.Command(path, args...)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("%s: %s\n", path, err)
	}
}

func main() {
	log.Println("Starting workbenchd...")
	// server := &qrpc.Server{}
	// l, err := mux.ListenTCP(busAddr)
	// if err != nil {
	// 	panic(err)
	// }
	// go func() {
	// 	fmt.Printf("Listening on %s...\n", busAddr)
	// 	log.Fatal(server.Serve(l, bus.NewBus()))
	// }()

	s, err := supervisor.NewSupervisor(os.Stdout)
	s.ChangeCallback = func(path string, reloadable bool, deleted bool) {
		log.Println("Modules folder changed...")
		if !deleted {
			triggerHook("change", path)
		}
	}
	if err != nil {
		panic(err)
	}
	go s.Watch()
	err = s.LoadDir("modules")
	if err != nil {
		panic(err)
	}
	s.Wait()
	log.Println("All modules finished.")
}
