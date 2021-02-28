//
// NOT USED!!!!
//
package swedbank

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/hultan/softfond/internal/data"
	"github.com/hultan/softfond/internal/tools"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	constUrl = "https://spara.swedbank.se/app/fondlista/fond/{FundId}?bankId=08999"
)

type Swedbank struct {
}

func NewSwedbank() *Swedbank {
	return new(Swedbank)
}

func (s *Swedbank) GetTodaysRate(fund *data.Fund) {
	// Request the HTML page.
	url := strings.Replace(constUrl, "{FundId}", fund.Identifier, 1)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(doc.Html())
	text := doc.Find("section._content-block").First().Text()
	fmt.Println(text)
	todaysRateString := tools.GetRateString(text)
	todaysRate, _ := strconv.ParseFloat(todaysRateString, 64)
	fund.LatestRate = todaysRate

	err = res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Swedbank) PrintFund(fund *data.Fund) {
	fmt.Printf("%.30s : %9.4f (%s) %s\n", fund.NameFormat(30), fund.LatestRate, fund.BuyingRateFormat(), fund.ProfitLossPercentFormat())
}

func (s *Swedbank) CalculateFundsTotal(funds *data.Funds) {
	for id := range funds.List {
		fund := funds.List[id]
		funds.TotalPurchasePrice += fund.PurchasePrice
		funds.TotalValue += fund.Shares * fund.LatestRate
	}
	funds.ProfitLossPercent = funds.TotalValue/funds.TotalPurchasePrice*100 - 100
}

func (s *Swedbank) PrintFundsTotal(funds *data.Funds) {
	fmt.Printf("PURCHASE PRICE : %v\tVALUE : %v\tPROFIT/LOSS : %v%%", funds.TotalPurchasePrice, funds.TotalValue, funds.ProfitLossPercent)
}
//
// PRIVATE FUNCTIONS
//

