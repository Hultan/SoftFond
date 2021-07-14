package softfond

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softfond/internal/data"
	"github.com/hultan/softfond/internal/tools"
	"github.com/hultan/softteam-tools/pkg/resources"
	"log"
	"os"
)

type MainForm struct {
	Window                  *gtk.ApplicationWindow
	builder                  *tools.SoftBuilder
	TreeView                *gtk.TreeView
	AboutDialog             *gtk.AboutDialog
	FundsValueLabel         *gtk.Label
	FundsPurchasePriceLabel *gtk.Label
	FundsProfitLossLabel    *gtk.Label
	ToolbarQuit             *gtk.ToolButton
	ToolbarRefresh          *gtk.ToolButton

	FundList *fundList
	Funds    *data.Funds
}

// MainFormNew : Creates a new MainForm object
func MainFormNew() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	m.builder = tools.NewSoftBuilder("main.glade")

	// Controls & signals
	m.getControls()
	m.hookUpSignals()

	// Menu
	m.setupMenu()

	// Funds
	m.loadFunds()
	m.FundList = fundListNew(m.Funds, m.TreeView, m)
	m.updateTotals(m.Funds)

	// Set up main window
	m.Window.SetApplication(app)
	m.Window.SetTitle(applicationTitle + " " + applicationVersion)

	// CleanUp
	m.builder = nil

	// Show the main window
	m.Window.ShowAll()
}

func (m *MainForm) getControls() {
	// Get the main window from the glade file
	m.Window = m.builder.GetObject("main_window").(*gtk.ApplicationWindow)

	// Status bar
	statusBar := m.builder.GetObject("main_window_status_bar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId(applicationTitle), applicationTitle+" "+applicationVersion+", "+applicationCopyRight)

	// Get the tree view
	m.TreeView = m.builder.GetObject("fund_treeview").(*gtk.TreeView)

	// Labels
	m.FundsPurchasePriceLabel = m.builder.GetObject("total_purchase_price_value").(*gtk.Label)
	m.FundsValueLabel = m.builder.GetObject("total_value").(*gtk.Label)
	m.FundsProfitLossLabel = m.builder.GetObject("total_profit_loss_value").(*gtk.Label)

	// Toolbar quit button
	m.ToolbarQuit = m.builder.GetObject("toolbar_quit").(*gtk.ToolButton)

	// Toolbar refresh button
	m.ToolbarRefresh = m.builder.GetObject("toolbar_refresh").(*gtk.ToolButton)
}

func (m *MainForm) hookUpSignals() {
	// Hook up the destroy event
	_, err := m.Window.Connect("destroy", m.shutDown)
	if err != nil {
		log.Println("Failed to connect the MainForm.destroy event")
		log.Fatal(err)
	}

	// Hook up the toolbar quit button clicked signal
	_, err = m.ToolbarQuit.Connect("clicked", m.shutDown)
	if err != nil {
		log.Println("Failed to connect the toolbar_quit.clicked event")
		log.Fatal(err)
	}

	// Hook up the toolbar refresh button clicked signal
	// For some reason this does not work, see : https://forum.golangbridge.org/t/reciever-nil-problem/22708/3
	//_, err = m.ToolbarRefresh.Connect("clicked", m.FundList.updateFundsValue)
	_, err = m.ToolbarRefresh.Connect("clicked", func() {
		m.FundList.updateFundsValue()
		m.updateTotals(m.Funds)

	})
	if err != nil {
		log.Println("Failed to connect the toolbar_refresh.clicked event")
		log.Fatal(err)
	}
}

func (m *MainForm) loadFunds() {
	funds := data.FundsNew()

	// Load
	err := funds.Load()
	if err != nil {
		log.Fatal(err)
	}

	m.Funds = funds
	m.Funds.CalculateFundsTotalValue()
}

func (m *MainForm) setupMenu() {
	menuQuit := m.builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	_, err := menuQuit.Connect("activate", m.shutDown)
	if err != nil {
		log.Println("failed to connect menu_file_quit.activate signal")
		log.Fatal(err)
	}

	menuHelpAbout := m.builder.GetObject("menu_help_about").(*gtk.MenuItem)
	if err != nil {
		log.Println("failed to find menu item menu_help_about")
		log.Fatal(err)
	}
	_, err = menuHelpAbout.Connect("activate", m.openAboutDialog)
	if err != nil {
		log.Println("failed to connect menu_help_about.activate signal")
		log.Fatal(err)
	}
}

func (m *MainForm) openAboutDialog() {
	if m.AboutDialog == nil {
		about := m.builder.GetObject("about_dialog").(*gtk.AboutDialog)
		about.SetDestroyWithParent(true)
		about.SetTransientFor(m.Window)
		about.SetProgramName(applicationTitle)
		about.SetComments(applicationDescription)
		about.SetVersion(applicationVersion)
		about.SetCopyright(applicationCopyRight)
		resource := resources.NewResources()
		image, err := gdk.PixbufNewFromFile(resource.GetResourcePath("application.png"))
		if err == nil {
			about.SetLogo(image)
		}
		about.SetModal(true)
		about.SetPosition(gtk.WIN_POS_CENTER)

		_, err = about.Connect("response", func(dialog *gtk.AboutDialog, responseId gtk.ResponseType) {
			if responseId == gtk.RESPONSE_CANCEL || responseId == gtk.RESPONSE_DELETE_EVENT {
				about.Hide()
			}
		})
		if err != nil {
			log.Println("failed to connect about_dialog.response signal")
			log.Fatal(err)
		}

		m.AboutDialog = about
	}

	m.AboutDialog.Present()
}

func (m *MainForm) updateTotals(funds *data.Funds) {
	m.FundsPurchasePriceLabel.SetText(funds.PurchasePriceFormat())
	m.FundsValueLabel.SetText(funds.ValueFormat())
	m.FundsProfitLossLabel.SetText(funds.ProfitLossFormat())
}

func (m *MainForm) shutDown() {
	if m.FundList != nil {
		m.FundList.Destroy()
		m.FundList = nil
	}

	if m.ToolbarQuit != nil {
		m.ToolbarQuit.Destroy()
	}
	if m.ToolbarRefresh != nil {
		m.ToolbarRefresh.Destroy()
	}
	if m.FundsProfitLossLabel != nil {
		m.FundsProfitLossLabel.Destroy()
	}
	if m.FundsPurchasePriceLabel != nil {
		m.FundsPurchasePriceLabel.Destroy()
	}
	if m.FundsValueLabel != nil {
		m.FundsValueLabel.Destroy()
	}
	if m.AboutDialog != nil {
		m.AboutDialog.Destroy()
	}
	if m.TreeView !=nil {
		m.TreeView.Destroy()
	}

	if m.Window!=nil {
		m.Window.Destroy()
	}
}
