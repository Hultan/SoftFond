package softfond

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softfond/internal/data"
	"github.com/hultan/softfond/internal/morningstar"
	"log"
)

type fundList struct {
	mainForm  *MainForm
	listStore *gtk.ListStore
	funds     *data.Funds
}

// fundListNew : Creates a new fundList struct
func fundListNew(mainForm *MainForm) *fundList {
	f := new(fundList)
	f.mainForm = mainForm

	f.loadFunds()

	return f
}

func (f *fundList) loadFunds() {
	funds := data.FundsNew()

	// Load
	err := funds.Load()
	if err != nil {
		log.Fatal(err)
	}

	f.funds = funds
}

func (f *fundList) updateFundsValue() {

	go func() {
		morningStar := morningstar.New()
		for _, fund := range f.funds.List {
			morningStar.GetFundRate(fund)
		}

		f.funds.Save()

		f.refreshFundList()
	}()

}

// refresh : Refreshes the video list
func (f *fundList) refreshFundList() {
	var err error

	if f.listStore != nil {
		f.listStore.Clear()
	}

	f.mainForm.TreeView.SetModel(nil)
	f.listStore, err = gtk.ListStoreNew(
		glib.TYPE_STRING, // Fund name
		glib.TYPE_STRING, // Fund value
		gdk.PixbufGetType(),
		glib.TYPE_STRING, // Profit/Loss percent
		gdk.PixbufGetType(),
		glib.TYPE_STRING, // Short Term Profit/Loss percent
		glib.TYPE_STRING, // Background color
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, fund := range f.funds.List {
		f.addFundToList(fund, f.listStore)
	}

	f.mainForm.TreeView.SetModel(f.listStore)
}

func (f *fundList) addFundToList(fund *data.Fund, listStore *gtk.ListStore) {
	// Append fund to list
	iter := listStore.Append()
	err := listStore.Set(iter, []int{columnName, columnValue, columnTrend, columnProfitLoss, columnShortTermTrend, columnShortTermProfitLoss, columnBackground},
		[]interface{}{
			f.getNameColumn(fund),
			f.getValueColumn(fund),
			f.getTrendImageColumn(fund, false),
			f.getProfitLossColumn(fund, false),
			f.getTrendImageColumn(fund, true),
			f.getProfitLossColumn(fund, true),
			"White",
		})

	if err != nil {
		log.Fatal(err)
	}
}

// setupColumns : Sets up the listview columns
func (f *fundList) setupColumns() {
	helper := new(treeviewHelper)
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Fondnamn", columnName, columnNameWidth))
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Värde", columnValue, columnValueWidth))
	f.mainForm.TreeView.AppendColumn(helper.createImageColumn("Lång", columnTrend, columnTrendWidth))
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Långtids V/F", columnProfitLoss, columnProfitLossWidth))
	f.mainForm.TreeView.AppendColumn(helper.createImageColumn("Kort", columnShortTermTrend, columnShortTermTrendWidth))
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Korttids V/F", columnShortTermProfitLoss, columnShortTermProfitLossWidth))
}

func (f *fundList) getTrendImageColumn(fund *data.Fund, shortTerm bool) *gdk.Pixbuf {
	var thumbnailPath string = "assets/trend_up.png"

	if shortTerm {
		if fund.ShortTermProfitLossPercent() == 0 {
			thumbnailPath = "assets/trend_none.png"
		} else if fund.ShortTermProfitLossPercent() < 0 {
			thumbnailPath = "assets/trend_down.png"
		}

	} else {
		if fund.ProfitLossPercent() == 0 {
			thumbnailPath = "assets/trend_none.png"
		} else if fund.ProfitLossPercent() < 0 {
			thumbnailPath = "assets/trend_down.png"
		}
	}

	thumbnail, err := gdk.PixbufNewFromFile(thumbnailPath)
	if err != nil {
		log.Fatal(err)
	}

	return thumbnail
}

func (f *fundList) getNameColumn(fund *data.Fund) string {
	return `<span font="Sans 16"><span foreground="#222222">` + fund.DisplayName + `</span></span>
<span font="Sans 12"><span foreground="#666666">` + fund.FundCompanyName + `</span></span>`
}

func (f *fundList) getValueColumn(fund *data.Fund) string {
	return `<span font="Sans 16"><span foreground="#222222">` + fund.CurrentValueFormat() + `</span></span>
<span font="Sans 12"><span foreground="#666666">(` + fund.PurchasePriceFormat() + `)</span></span>`
}

func (f *fundList) getProfitLossColumn(fund *data.Fund, shortTerm bool) string {
	if shortTerm {
		if fund.ShortTermProfitLossPercent() == 0 {
			return `<span font="Sans 14"><span foreground="#000000">` + fund.ShortTermProfitLossPercentFormat() + `</span></span>`
		} else if fund.ShortTermProfitLossPercent() >= 0 {
			return `<span font="Sans 14"><span foreground="#00FF00">` + fund.ShortTermProfitLossPercentFormat() + `</span></span>`
		}
		return `<span font="Sans 14"><span foreground="#FF0000">` + fund.ShortTermProfitLossPercentFormat() + `</span></span>`
	} else {
		if fund.ProfitLossPercent() == 0 {
			return `<span font="Sans 14"><span foreground="#000000">` + fund.ProfitLossPercentFormat() + `</span></span>`
		} else if fund.ProfitLossPercent() >= 0 {
			return `<span font="Sans 14"><span foreground="#00FF00">` + fund.ProfitLossPercentFormat() + `</span></span>`
		}
		return `<span font="Sans 14"><span foreground="#FF0000">` + fund.ProfitLossPercentFormat() + `</span></span>`
	}
}
