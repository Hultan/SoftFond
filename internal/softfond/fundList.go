package softfond

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softfond/internal/data"
	"log"
)

type fundList struct {
	mainForm  *MainForm
	listStore *gtk.ListStore
}

// fundListNew : Creates a new fundList struct
func fundListNew(mainForm *MainForm) *fundList {
	f := new(fundList)
	f.mainForm = mainForm
	return f
}

func (f *fundList) refreshFundList() {

	funds := data.FundsNew()

	// Load
	err := funds.Load()
	if err != nil {
		log.Fatal(err)
	}

	f.refresh(funds)
}

// refresh : Refreshes the video list
func (f *fundList) refresh(funds *data.Funds) {
	var err error

	if f.listStore != nil {
		f.listStore.Clear()
	}

	f.mainForm.TreeView.SetModel(nil)
	f.listStore, err = gtk.ListStoreNew(
		gdk.PixbufGetType(),
		glib.TYPE_STRING, // Fund name
		glib.TYPE_STRING, // Fund value
		glib.TYPE_STRING, // Profit/Loss percent
		glib.TYPE_STRING, // Background color
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, fund := range funds.List {
		f.addFundToList(fund, f.listStore)
	}

	f.mainForm.TreeView.SetModel(f.listStore)
}

func (f *fundList) addFundToList(fund *data.Fund, listStore *gtk.ListStore) {
	// Append fund to list
	iter := listStore.Append()
	err := listStore.Set(iter, []int{columnTrend, columnName, columnValue, columnProfitLoss, columnBackground},
		[]interface{}{
			f.getTrendImageColumn(fund),
			f.getNameColumn(fund),
			f.getValueColumn(fund),
			f.getProfitLossColumn(fund),
			"White",
		})

	if err != nil {
		log.Fatal(err)
	}
}

// setupColumns : Sets up the listview columns
func (f *fundList) setupColumns() {
	helper := new(treeviewHelper)
	f.mainForm.TreeView.AppendColumn(helper.createImageColumn("Trend", columnTrend, columnTrendWidth))
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Fondnamn", columnName, columnNameWidth))
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Värde", columnValue, columnValueWidth))
	f.mainForm.TreeView.AppendColumn(helper.createTextColumn("Procent", columnProfitLoss, columnProfitLossWidth))
}

func (f *fundList) getTrendImageColumn(fund *data.Fund) *gdk.Pixbuf {
	var thumbnailPath string
	if fund.ProfitLossPercent()>=0 {
		thumbnailPath = "assets/trend_up.png"
	} else {
		thumbnailPath = "assets/trend_down.png"
	}

	thumbnail, err := gdk.PixbufNewFromFile(thumbnailPath)
	if err != nil {
		log.Fatal(err)
	}

	return thumbnail
}

func (f *fundList) getNameColumn(fund *data.Fund) string {
	return `<span font="Sans 16"><span foreground="#222222">` + fund.DisplayName + `</span></span>
<span font="Sans 12"><span foreground="#666666">` + fund.FundCompany +`</span></span>`
}

func (f *fundList) getValueColumn(fund *data.Fund) string {
	return `<span font="Sans 16"><span foreground="#222222">` + fund.CurrentValueFormat() + `</span></span>
<span font="Sans 12"><span foreground="#666666">(` + fund.PurchasePriceFormat() + `)</span></span>`
}

func (f *fundList) getProfitLossColumn(fund *data.Fund) string {
	if fund.ProfitLossPercent()>=0 {
		return `<span font="Sans 14"><span foreground="#00FF00">` + fund.ProfitLossPercentFormat() + `</span></span>`
	}
	return `<span font="Sans 14"><span foreground="#FF0000">` + fund.ProfitLossPercentFormat() + `</span></span>`
}