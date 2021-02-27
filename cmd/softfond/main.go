package main

import (
	"github.com/hultan/softfond/internal/data"
	"github.com/hultan/softfond/internal/morningstar"
	"log"
)

func main() {
	funds := data.NewFunds()

	// Load
	err := funds.Load()
	if err != nil {
		log.Fatal(err)
	}

	morningstar := morningstar.NewMorningStar()
	for _,fund := range funds.List {
		err = morningstar.GetFundValue(fund)
		if err!=nil {
			log.Fatal(err)
		}
		morningstar.PrintFund(fund)
	}
	morningstar.GetFundsValue(funds)
	morningstar.PrintFunds(funds)

	// Save
	err = funds.Save()
	if err!=nil {
		log.Fatal(err)
	}
}

