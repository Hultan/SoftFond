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
