package tools

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type SoftBuilder struct {
	builder *gtk.Builder
}

func NewSoftBuilder(fileName string) *SoftBuilder {
	builder := new(SoftBuilder)
	builder.createBuilder(fileName)
	return builder
}

func (s *SoftBuilder) createBuilder(gladeFileName string) {
	gladePath, err := GetAssetsPath(gladeFileName)
	if err != nil {
		panic(err)
	}

	builder, err := gtk.BuilderNewFromFile(gladePath)
	if err != nil {
		panic(err)
	}

	s.builder = builder
}

func (s *SoftBuilder) GetObject(name string) glib.IObject {
	obj, err := s.builder.GetObject(name)
	if err != nil {
		panic(err)
	}

	return obj
}

