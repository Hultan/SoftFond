package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softfond/internal/softfond"
	"log"
	"os"
)

const (
	ApplicationId    = "se.softteam.softfond"
	ApplicationFlags = glib.APPLICATION_FLAGS_NONE
)

func main() {
	// Create an gtk.application
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		log.Fatal(err)
	}

	// Create the main form and hook up the activate signal for the application
	mainForm := softfond.MainFormNew()
	_, err = application.Connect("activate", mainForm.OpenMainForm)
	if err != nil {
		log.Fatal(err)
	}

	// Run the application
	os.Exit(application.Run(nil))
}
