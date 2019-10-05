package persistence

import (
	"encoding/json"
	"github.com/21stio/go-playground/gifts/pkg/types"
	"github.com/spf13/afero"
)

type JsonLoader struct {
	fs afero.Fs
	employeesPath string
	giftsPath string
}

func NewJsonLoader(fs afero.Fs, employeesPath, giftsPath string) (l JsonLoader) {
	l.fs = fs
	l.employeesPath = employeesPath
	l.giftsPath = giftsPath

	return
}

func (l JsonLoader) LoadEmployees() (es []types.Employee, err error) {
	f, err := l.fs.Open(l.employeesPath)
	if err != nil {
	    return
	}

	err = json.NewDecoder(f).Decode(&es)
	if err != nil {
	    return
	}

	return
}

func (l JsonLoader) LoadGifts() (gs []types.Gift, err error) {
	f, err := l.fs.Open(l.giftsPath)
	if err != nil {
		return
	}

	err = json.NewDecoder(f).Decode(&gs)
	if err != nil {
		return
	}

	return
}