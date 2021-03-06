package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Funds struct {
	List                   []*Fund `json:"Funds"`
	TotalPurchasePrice     float64
	TotalValue             float64
	TotalProfitLossPercent float64
}

// FundsNew : Create a new funds struct
func FundsNew() *Funds {
	f := new(Funds)
	return f
}

// Load : Load funds from a json file
func (f *Funds) Load() error {
	bytes, err := ioutil.ReadFile("config/funds.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, f)
	if err != nil {
		return err
	}

	return nil
}

// Save : Save funds to a json file
func (f *Funds) Save() error {
	bytes, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config/funds.json", bytes, 0644)
	if err != nil {
		return err
	}
	return err
}

func (f *Funds) CalculateFundsTotalValue() {
	f.TotalPurchasePrice = 0
	f.TotalValue = 0

	for id := range f.List {
		fund := f.List[id]
		f.TotalPurchasePrice += fund.PurchasePrice
		f.TotalValue += fund.Shares * fund.TodaysRate
	}
	f.TotalProfitLossPercent = f.TotalValue/f.TotalPurchasePrice*100 - 100
}

func (f *Funds) PurchasePriceFormat() string {
	return fmt.Sprintf("%.0f SEK", f.TotalPurchasePrice)
}

func (f *Funds) ValueFormat() string {
	return fmt.Sprintf("%.0f SEK", f.TotalValue)
}

func (f *Funds) ProfitLossFormat() string {
	return fmt.Sprintf("%6.2f %%", f.TotalProfitLossPercent)
}

