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
	application, err := gtk.ApplicationNew(ApplicationId, ApplicationFlags)
	if err != nil {
		log.Fatal(err)
	}

	mainForm := softfond.MainFormNew()
	_, err = application.Connect("activate", mainForm.OpenMainForm)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(application.Run(nil))
}

//func test() {
//	funds := data.NewFunds()
//
//	// Load
//	err := funds.Load()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	morningstar := morningstar.NewMorningStar()
//	for _, fund := range funds.List {
//		err = morningstar.GetFundValue(fund)
//		if err != nil {
//			log.Fatal(err)
//		}
//		morningstar.PrintFund(fund)
//	}
//	morningstar.GetFundsValue(funds)
//	morningstar.PrintFunds(funds)
//
//	// Save
//	err = funds.Save()
//	if err != nil {
//		log.Fatal(err)
//	}
//}
