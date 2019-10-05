package boot

import (
	"fmt"
	"github.com/21stio/go-playground/gifts/pkg/business"
	"github.com/21stio/go-playground/gifts/pkg/communication"
	"github.com/21stio/go-playground/gifts/pkg/persistence"
	"github.com/21stio/go-playground/gifts/pkg/types"
	"github.com/spf13/afero"
	"log"
)

func GetDependencies(employeesPath, giftsPath, addr string) (d types.Dependencies, err error) {
	fs := afero.NewOsFs()

	l := persistence.NewJsonLoader(fs, employeesPath, giftsPath)

	es, err := l.LoadEmployees()
	if err != nil {
		return
	}
	requireUniqueEmployeeNames(es)

	gs, err := l.LoadGifts()
	if err != nil {
		return
	}
	requireUniqueGiftNames(gs)

	d.Storage = persistence.NewMemoryStorage(es, gs)

	logErrFunc := func(args ...interface{}) { log.Println(append([]interface{}{"err:"}, args...)...) }
	facade := business.NewFacade(d.Storage, logErrFunc, log.Println)

	d.Server = communication.NewServer(facade, addr, logErrFunc, log.Println)

	return
}

func requireUniqueEmployeeNames(es []types.Employee) {
	m := map[string]bool{}
	for _, e := range es {
		_, ok := m[e.Name]
		if ok {
			panic(fmt.Sprintf("employee names must be unique, name %v occurs at least twice", e.Name))
		}

		m[e.Name] = true
	}
}

func requireUniqueGiftNames(gs []types.Gift) {
	m := map[string]bool{}
	for _, g := range gs {
		_, ok := m[g.Name]
		if ok {
			panic(fmt.Sprintf("gift names must be unique, name %v occurs at least twice", g.Name))
		}

		m[g.Name] = true
	}
}