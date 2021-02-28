package data

import (
	"encoding/json"
	"fmt"
	"github.com/hultan/softfond/internal/tools"
	"io/ioutil"
)

type Funds struct {
	List               []*Fund `json:"Funds"`
	TotalPurchasePrice float64
	TotalValue         float64
	ProfitLossPercent  float64
}

type Fund struct {
	Id            int     `json:"Id"`
	Name          string  `json:"Name"`
	DisplayName   string  `json:"DisplayName"`
	FundCompany   string  `json:"FundCompany"`
	ParserName    string  `json:"ParserName"`
	LatestRate    float64 `json:"LatestRate"`
	Shares        float64 `json:"Shares"`
	PurchasePrice float64 `json:"PurchasePrice"`
	Identifier    string  `json:"Identifier"`
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

func (f *Fund) BuyingRate() float64 {
	return f.PurchasePrice / f.Shares
}

func (f *Fund) ProfitLossPercent() float64 {
	return f.LatestRate/f.BuyingRate()*100 - 100
}

func (f *Fund) NameFormat(maxLength int) string {
	name := []rune(f.Name)
	l := maxLength - len(name)
	if l > 0 {
		return f.Name + tools.MultiplySpaces(l-1)
	}
	return string(name[0:maxLength])
}

func (f *Fund) LatestRateFormat() string {
	return fmt.Sprintf("%9.4f", f.LatestRate)
}

func (f *Fund) BuyingRateFormat() string {
	return fmt.Sprintf("%9.4f", f.BuyingRate())
}

func (f *Fund) PurchasePriceFormat() string {
	return fmt.Sprintf("%.0f SEK", f.PurchasePrice)
}

func (f *Fund) CurrentValueFormat() string {
	return fmt.Sprintf("%.0f SEK", f.LatestRate * f.Shares)
}

func (f *Fund) ProfitLossPercentFormat() string {
	return fmt.Sprintf("%6.2f%%", f.ProfitLossPercent())
}
