package softfond

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softfond/internal/data"
	"github.com/hultan/softfond/internal/tools"
	gtkHelper "github.com/hultan/softteam-tools/pkg/gtk-helper"
	"github.com/hultan/softteam-tools/pkg/resources"
	"log"
	"os"
)

type MainForm struct {
	Window                  *gtk.ApplicationWindow
	Helper                  *gtkHelper.GtkHelper
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

	// Create a new gtk helper
	m.Helper = m.createHelper()

	// Controls & signals
	m.getControls()
	m.hookUpSignals()

	// Menu
	m.setupMenu()

	// Funds
	m.loadFunds()
	m.FundList = fundListNew(m.Funds, m.TreeView)
	m.updateTotals(m.Funds)

	// Set up main window
	m.Window.SetApplication(app)
	m.Window.SetTitle(applicationTitle + " " + applicationVersion)

	// CleanUp
	m.Helper = nil

	// Show the main window
	m.Window.ShowAll()
}

func (m *MainForm) createHelper() *gtkHelper.GtkHelper {
	return gtkHelper.GtkHelperNew(m.createGuilder())
}

func (m *MainForm) createGuilder() *gtk.Builder {
	builder, err := gtk.BuilderNewFromFile(tools.GetResourcePath("assets/main.glade"))
	if err != nil {
		log.Println("Failed to create builder")
		log.Fatal(err)
	}
	return builder
}

func (m *MainForm) getControls() {
	// Get the main window from the glade file
	window, err := m.Helper.GetApplicationWindow("main_window")
	if err != nil {
		log.Println("Failed to find main_window")
		log.Fatal(err)
	}
	m.Window = window

	// Status bar
	statusBar, err := m.Helper.GetStatusBar("main_window_status_bar")
	if err != nil {
		log.Println("Failed to find main_window_status_bar")
		log.Fatal(err)
	}
	statusBar.Push(statusBar.GetContextId(applicationTitle), applicationTitle+" "+applicationVersion+", "+applicationCopyRight)

	// Get the tree view
	treeView, err := m.Helper.GetTreeView("fund_treeview")
	if err != nil {
		log.Fatal(err)
	}
	m.TreeView = treeView

	// Labels
	label, err := m.Helper.GetLabel("total_purchase_price_value")
	if err != nil {
		log.Println("Failed to find total_purchase_price_value")
		log.Fatal(err)
	}
	m.FundsPurchasePriceLabel = label
	label, err = m.Helper.GetLabel("total_value")
	if err != nil {
		log.Println("Failed to find total_value")
		log.Fatal(err)
	}
	m.FundsValueLabel = label
	label, err = m.Helper.GetLabel("total_profit_loss_value")
	if err != nil {
		log.Println("Failed to find total_profit_loss_value")
		log.Fatal(err)
	}
	m.FundsProfitLossLabel = label

	// Toolbar quit button
	button, err := m.Helper.GetToolButton("toolbar_quit")
	if err != nil {
		log.Println("Failed to find toolbar_quit")
		log.Fatal(err)
	}
	m.ToolbarQuit = button

	// Toolbar refresh button
	button, err = m.Helper.GetToolButton("toolbar_refresh")
	if err != nil {
		log.Println("Failed to find toolbar_refresh")
		log.Fatal(err)
	}
	m.ToolbarRefresh = button
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
	//_, err = m.ToolbarRefresh.Connect("clicked", m.FundList.updateFundsValue)
	_, err = m.ToolbarRefresh.Connect("clicked", func() {
		m.FundList.updateFundsValue()
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
	menuQuit, err := m.Helper.GetMenuItem("menu_file_quit")
	if err != nil {
		log.Println("failed to find menu item menu_file_quit")
		log.Fatal(err)
	}
	_, err = menuQuit.Connect("activate", m.shutDown)
	if err != nil {
		log.Println("failed to connect menu_file_quit.activate signal")
		log.Fatal(err)
	}

	menuHelpAbout, err := m.Helper.GetMenuItem("menu_help_about")
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
		about, err := m.Helper.GetAboutDialog("about_dialog")
		if err != nil {
			log.Println("failed to find dialog about_dialog")
			log.Fatal(err)
		}
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
