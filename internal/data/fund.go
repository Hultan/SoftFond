package data

import (
	"fmt"
	"github.com/hultan/softfond/internal/tools"
	"time"
)

type Fund struct {
	Id              int     `json:"Id"`
	FundName        string  `json:"FundName"`
	DisplayName     string  `json:"DisplayName"`
	FundCompanyName string  `json:"FundCompanyName"`
	ParserName      string  `json:"ParserName"`
	FundIdentifier  string  `json:"FundIdentifier"`
	Shares          float64 `json:"Shares"`
	PurchasePrice   float64 `json:"PurchasePrice"`

	TodaysRate           float64   `json:"TodaysRate"`
	TodaysUpdateTime     time.Time `json:"TodaysUpdateTime"`
	YesterdaysRate       float64   `json:"YesterdaysRate"`
	YesterdaysUpdateTime time.Time `json:"YesterdaysUpdateTime"`
}

func (f *Fund) BuyingRate() float64 {
	return f.PurchasePrice / f.Shares
}

func (f *Fund) ProfitLossPercent() float64 {
	return f.TodaysRate/f.BuyingRate()*100 - 100
}

func (f *Fund) ShortTermProfitLossPercent() float64 {
	return f.TodaysRate/f.YesterdaysRate*100-100
}

func (f *Fund) NameFormat(maxLength int) string {
	name := []rune(f.FundName)
	l := maxLength - len(name)
	if l > 0 {
		return f.FundName + tools.MultiplySpaces(l-1)
	}
	return string(name[0:maxLength])
}

func (f *Fund) CurrentRateFormat() string {
	return fmt.Sprintf("%9.4f", f.TodaysRate)
}

func (f *Fund) BuyingRateFormat() string {
	return fmt.Sprintf("%9.4f", f.BuyingRate())
}

func (f *Fund) PurchasePriceFormat() string {
	return fmt.Sprintf("%.0f SEK", f.PurchasePrice)
}

func (f *Fund) CurrentValueFormat() string {
	return fmt.Sprintf("%.0f SEK", f.TodaysRate*f.Shares)
}

func (f *Fund) ProfitLossPercentFormat() string {
	return fmt.Sprintf("%6.2f %%", f.ProfitLossPercent())
}

func (f *Fund) ShortTermProfitLossPercentFormat() string {
	return fmt.Sprintf("%6.2f %%", f.ShortTermProfitLossPercent())
}


