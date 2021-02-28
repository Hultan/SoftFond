package softfond

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/hultan/softfond/internal/tools"
	gtkHelper "github.com/hultan/softteam-tools/pkg/gtk-helper"
	"github.com/hultan/softteam-tools/pkg/resources"
	"log"
	"os"
)

type MainForm struct {
	Window      *gtk.ApplicationWindow
	Helper      *gtkHelper.GtkHelper
	TreeView    *gtk.TreeView
	AboutDialog *gtk.AboutDialog
	FundList    *fundList
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
	builder, err := gtk.BuilderNewFromFile(tools.GetResourcePath("../assets", "main.glade"))
	if err != nil {
		log.Println("Failed to create builder")
		log.Fatal(err)
	}
	helper := gtkHelper.GtkHelperNew(builder)
	m.Helper = helper

	// Get the main window from the glade file
	window, err := helper.GetApplicationWindow("main_window")
	if err != nil {
		log.Println("Failed to find main_window")
		log.Fatal(err)
	}
	m.Window = window

	// Set up main window
	window.SetApplication(app)
	window.SetTitle(applicationTitle + " " + applicationVersion)

	// Hook up the destroy event
	_, err = window.Connect("destroy", window.Close)
	if err != nil {
		log.Println("Failed to connect the mainForm.destroy event")
		log.Fatal(err)
	}

	// Quit button
	button, err := helper.GetToolButton("toolbar_quit")
	if err != nil {
		log.Println("Failed to find toolbar_quit")
		log.Fatal(err)
	}
	_, err = button.Connect("clicked", window.Close)
	if err != nil {
		log.Println("Failed to connect the toolbar_quit.clicked event")
		log.Fatal(err)
	}

	// Status bar
	statusBar, err := helper.GetStatusBar("main_window_status_bar")
	if err != nil {
		log.Println("Failed to find main_window_status_bar")
		log.Fatal(err)
	}
	statusBar.Push(statusBar.GetContextId(applicationTitle),  applicationTitle + " " + applicationVersion + ", " + applicationCopyRight)

	// Menu
	m.setupMenu(window)

	// Get the tree view
	treeView, err := helper.GetTreeView("fund_treeview")
	if err != nil {
		log.Fatal(err)
	}
	m.TreeView = treeView

	// Setup fund list
	m.FundList = fundListNew(m)
	m.FundList.setupColumns()
	m.FundList.refreshFundList()

	// Refresh button
	button, err = helper.GetToolButton("toolbar_refresh")
	if err != nil {
		log.Println("Failed to find toolbar_refresh")
		log.Fatal(err)
	}
	_, err = button.Connect("clicked", m.FundList.updateFundsValue)
	if err != nil {
		log.Println("Failed to connect the toolbar_refresh.clicked event")
		log.Fatal(err)
	}

	// Show the main window
	window.ShowAll()
}

func (m *MainForm) setupMenu(window *gtk.ApplicationWindow) {
	menuQuit, err := m.Helper.GetMenuItem("menu_file_quit")
	if err != nil {
		log.Println("failed to find menu item menu_file_quit")
		log.Fatal(err)
	}
	_, err = menuQuit.Connect("activate", window.Close)
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
