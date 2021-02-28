package softfond

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

// TreeviewHelper : A helper class for a gtk treeview
type TreeviewHelper struct {
}

// CreateTextColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *TreeviewHelper) CreateTextColumn(title string, id int, width int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "markup", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.AddAttribute(cellRenderer, "background", ColumnBackground)

	if width<0 {
		column.SetExpand(true)
	} else {
		column.SetFixedWidth(width)
	}

	return column
}

// CreateImageColumn : Add a column to the tree view (during the initialization of the tree view)
func (t *TreeviewHelper) CreateImageColumn(title string, imageColumn int, width int) *gtk.TreeViewColumn {
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
