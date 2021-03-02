package morningstar

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hultan/softfond/internal/data"
	"github.com/hultan/softfond/internal/tools"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	constUrl = "https://www.morningstar.se/se/funds/snapshot/snapshot.aspx?id={FundId}"
)

type Morningstar struct {
}

// New : Creates a new MorningStar struct
func New() *Morningstar {
	return new(Morningstar)
}

// GetFundRate : Gets todays fund rate from MorningStar.se
func (m *Morningstar) GetFundRate(fund *data.Fund) error {
	// Request the HTML page.
	url := strings.Replace(constUrl, "{FundId}", fund.FundIdentifier, 1)
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("status code error: %d %s", res.StatusCode, res.Status))
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	text := doc.Find("table.overviewKeyStatsTable tr:nth-child(2) td.text").First().Text()
	todaysRateString := tools.GetRateString(text)
	todaysRate, err := strconv.ParseFloat(todaysRateString, 64)
	if err != nil {
		return err
	}

	// If enough time has passed, move today to yesterday
	if time.Now().YearDay() > fund.TodaysUpdateTime.YearDay() {
		fund.YesterdaysRate = fund.TodaysRate
		fund.YesterdaysUpdateTime = fund.TodaysUpdateTime
	}

	// Set todays rate and update time
	fund.TodaysRate = todaysRate
	fund.TodaysUpdateTime = time.Now()

	err = res.Body.Close()
	if err != nil {
		return err
	}

	fmt.Printf("Fonden '%s' är uppdaterad, dagskurs : %v...\n", fund.FundName, fund.TodaysRate)

	return nil
}

func (m *Morningstar) PrintFund(fund *data.Fund) {
	fmt.Printf("%.30s : %9.4f (%s) %s\n", fund.NameFormat(30), fund.TodaysRate, fund.BuyingRateFormat(), fund.ProfitLossPercentFormat())
}

// Todo : Move out
func (m *Morningstar) GetFundsValue(funds *data.Funds) {
	funds.TotalPurchasePrice = 0
	funds.TotalValue = 0

	for id := range funds.List {
		fund := funds.List[id]
		funds.TotalPurchasePrice += fund.PurchasePrice
		funds.TotalValue += fund.Shares * fund.TodaysRate
	}
	funds.ProfitLossPercent = funds.TotalValue/funds.TotalPurchasePrice*100 - 100
}

func (m *Morningstar) PrintFunds(funds *data.Funds) {
	fmt.Printf("PURCHASE PRICE : %v\tVALUE : %v\tPROFIT/LOSS : %v%%", funds.TotalPurchasePrice, funds.TotalValue, funds.ProfitLossPercent)
}
