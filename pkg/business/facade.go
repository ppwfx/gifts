package business

import (
	"fmt"
	"github.com/21stio/go-playground/gifts/pkg/types"
	"github.com/pkg/errors"
	"sync"
)

type Facade struct {
	storage      types.Storage
	logErrorFunc func(...interface{})
	logInfoFunc  func(...interface{})
	mux          *sync.Mutex
}

func NewFacade(storage types.Storage, logErrorFunc func(...interface{}), logInfoFunc func(...interface{})) (*Facade) {
	return &Facade{
		storage:      storage,
		logErrorFunc: logErrorFunc,
		logInfoFunc:  logInfoFunc,
		mux:          &sync.Mutex{},
	}
}

func (f Facade) AssignEmployeeToGift(employeeName string) (err error) {
	f.mux.Lock()
	defer f.mux.Unlock()

	e, err := f.storage.GetEmployeeByName(employeeName)
	if err != nil {
		return
	}

	if e.AssignedToGift != "" {
		return errors.Errorf("employee is already assigned to a gift")
	}

	es, err := f.storage.GetUnassignedEmployees()
	if err != nil {
		return
	}

	gs, err := f.storage.GetUnassignedGifts()
	if err != nil {
		return
	}

	employeesToGiftsToDemands := getEmployeesToGiftsToDemands(es, gs)

	giftsToDemands, ok := employeesToGiftsToDemands[employeeName]
	if !ok {
		return errors.Errorf("¯\\_(ツ)_/¯")
	}

	giftName := match(giftsToDemands)
	if giftName == "" {
		return errors.Errorf("¯\\_(ツ)_/¯")
	}

	err = f.storage.AssignEmployeeToGift(employeeName, giftName)
	if err != nil {
		return
	}

	f.logInfoFunc(fmt.Sprintf("assigned employee: %v to gift: %v", employeeName, giftName))

	return
}

func match(giftsToDemand map[string]Demand) (g string) {
	if len(giftsToDemand) == 0 {
		return ""
	}

	if len(giftsToDemand) == 1 {
		for k, _ := range giftsToDemand {
			return k
		}
	}

	var greatestMinRange int
	var greatestMinRangeGift string
	var greatestMinRangeTotal int
	for g, d := range giftsToDemand {
		greatestMinRangeGift = g
		greatestMinRangeTotal = d.Total
		greatestMinRange = min(d.Range)
	}

	for g, d := range giftsToDemand {
		minRange := min(d.Range)
		if minRange > greatestMinRange {
			greatestMinRangeGift = g
			greatestMinRangeTotal = d.Total
			greatestMinRange = minRange
		}
		if minRange == greatestMinRange &&
			d.Total < greatestMinRangeTotal {
			greatestMinRangeGift = g
			greatestMinRangeTotal = d.Total
			greatestMinRange = minRange
		}
	}

	return greatestMinRangeGift
}

func min(numbers map[int]int) int {
	var minNumber int
	for minNumber = range numbers {
		break
	}
	for n := range numbers {
		if n < minNumber {
			minNumber = n
		}
	}
	return minNumber
}

type Demand struct {
	Total int
	Range map[int]int
}

func getEmployeesToGiftsToDemands(es []types.Employee, gs []types.Gift) (employeesToGiftsToDemands map[string]map[string]Demand) {
	employeesToGiftsToDemands = map[string]map[string]Demand{}

	interestsToGifts := getInterestsToGifts(gs)

	interestsToEmployees := getInterestsToEmployees(es)

	employeesToGifts := getEmployeesToGifts(es, interestsToGifts)

	giftsToEmployees := getGiftsToEmployees(gs, interestsToEmployees)

	giftsToDemands := map[string]Demand{}
	for _, g := range gs {
		giftsToDemands[g.Name] = Demand{
			Total: len(giftsToEmployees[g.Name]),
			Range: map[int]int{},
		}
	}

	for _, e := range es {
		for _, g := range employeesToGifts[e.Name] {
			giftsToDemands[g].Range[len(employeesToGifts[e.Name])] += 1
		}
	}

	for e, gs := range employeesToGifts {
		employeesToGiftsToDemands[e] = map[string]Demand{}
		for _, g := range gs {
			employeesToGiftsToDemands[e][g] = giftsToDemands[g]
		}
	}

	return
}

func getInterestsToGifts(gs []types.Gift) (interestsToGifts map[string][]string) {
	interestsToGifts = map[string][]string{}

	for _, g := range gs {
		for _, c := range g.Categories {
			interestsToGifts[c] = append(interestsToGifts[c], g.Name)
		}
	}

	return
}

func getInterestsToEmployees(es []types.Employee) (interestsToEmployees map[string][]string) {
	interestsToEmployees = map[string][]string{}
	for _, e := range es {
		for _, i := range e.Interests {
			interestsToEmployees[i] = append(interestsToEmployees[i], e.Name)
		}
	}

	return
}

func getEmployeesToGifts(es []types.Employee, interestsToGifts map[string][]string) (employeeToGifts map[string][]string) {
	employeeToGifts = map[string][]string{}

	employeeToGiftsMap := map[string]map[string]bool{}
	for _, e := range es {
		employeeToGiftsMap[e.Name] = map[string]bool{}
		for _, i := range e.Interests {
			for _, g := range interestsToGifts[i] {
				employeeToGiftsMap[e.Name][g] = true
			}
		}
	}

	for e, gs := range employeeToGiftsMap {
		for g, _ := range gs {
			employeeToGifts[e] = append(employeeToGifts[e], g)
		}
	}

	return
}

func getGiftsToEmployees(gs []types.Gift, interestsToEmployees map[string][]string) (giftsToEmployees map[string][]string) {
	giftsToEmployeesMap := map[string]map[string]bool{}
	for _, g := range gs {
		giftsToEmployeesMap[g.Name] = map[string]bool{}
		for _, c := range g.Categories {
			for _, e := range interestsToEmployees[c] {
				giftsToEmployeesMap[g.Name][e] = true
			}
		}
	}

	giftsToEmployees = map[string][]string{}
	for g, es := range giftsToEmployeesMap {
		for e, _ := range es {
			giftsToEmployees[g] = append(giftsToEmployees[g], e)
		}
	}

	return
}
