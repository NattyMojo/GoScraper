package main

import (
	"fmt"
	"log"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

func startDash() {
	// Set logger
	l := log.New(log.Writer(), log.Prefix(), log.Flags())

	// Create astilectron
	a, err := astilectron.New(l, astilectron.Options{
		AppName:           "GoScraping",
		BaseDirectoryPath: "src",
	})
	if err != nil {
		l.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

  // Add a listener on Astilectron
  a.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
      log.Println("App has crashed")
      return
  })

	// Start
	if err = a.Start(); err != nil {
		l.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// New window
	var w *astilectron.Window
	if w, err = a.NewWindow("index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(850),
		Width:  astikit.IntPtr(1200),
	}); err != nil {
		l.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

  // Add a listener on the window
  // This one alerts when the Window is resized
  // w.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
  //     log.Println("Window resized")
  //     return
  // })

	// Create windows
	if err = w.Create(); err != nil {
		l.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// Blocking pattern
	a.Wait()
}
