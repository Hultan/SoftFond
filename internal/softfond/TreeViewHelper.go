package softfond

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// treeviewHelper : A helper class for a gtk treeview
type treeviewHelper struct {
}

// createTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *treeviewHelper) createTextColumn(title string, id int, width int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "markup", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.AddAttribute(cellRenderer, "background", columnBackground)

	if width<0 {
		column.SetExpand(true)
	} else {
		column.SetFixedWidth(width)
	}

	return column
}

// createImageColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *treeviewHelper) createImageColumn(title string, imageColumn int, width int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererPixbufNew()
	if err != nil {
		log.Fatal("Unable to create pixbuf cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "pixbuf", imageColumn)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetFixedWidth(width)

	return column
}
