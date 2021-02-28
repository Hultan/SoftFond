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

const (
	ColumnTrend = iota
	ColumnName
	ColumnValue
	ColumnBackground
)

const (
	ColumnTrendWidth = 64
	ColumnNameWidth  = -1
	ColumnValueWidth = 200
)

func FundListNew(mainForm *MainForm) *fundList {
	f := new(fundList)
	f.mainForm = mainForm
	return f
}

// Refresh : Refreshes the video list
func (f *fundList) Refresh() {
	var err error

	if f.listStore != nil {
		f.listStore.Clear()
	}

	f.mainForm.TreeView.SetModel(nil)
	f.listStore, err = gtk.ListStoreNew(
		gdk.PixbufGetType(),
		glib.TYPE_STRING, // Fund name
		glib.TYPE_STRING, // Fund value
		glib.TYPE_STRING, // Background color
	)
	if err != nil {
		log.Fatal(err)
	}

	funds := data.NewFunds()

	// Load
	err = funds.Load()
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
	err := listStore.Set(iter, []int{ColumnTrend, ColumnName, ColumnValue, ColumnBackground},
		[]interface{}{
			f.getTrendImageColumn(fund),
			f.getNameColumn(fund),
			f.getValueColumn(fund),
			"White",
		})

	if err != nil {
		log.Fatal(err)
	}
}

// SetupColumns : Sets up the listview columns
func (f *fundList) SetupColumns() {
	helper := new(TreeviewHelper)
	f.mainForm.TreeView.AppendColumn(helper.CreateImageColumn("Trend", ColumnTrend, ColumnTrendWidth))
	f.mainForm.TreeView.AppendColumn(helper.CreateTextColumn("Fondnamn", ColumnName, ColumnNameWidth))
	f.mainForm.TreeView.AppendColumn(helper.CreateTextColumn("Värde", ColumnValue, ColumnValueWidth))
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