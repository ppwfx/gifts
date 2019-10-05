package persistence

import (
	"github.com/21stio/go-playground/gifts/pkg/types"
	"github.com/pkg/errors"
)

type MemoryStorage struct {
	employees map[string]types.Employee
	gifts     map[string]types.Gift
}

func NewMemoryStorage(employees []types.Employee, gifts []types.Gift) (m *MemoryStorage) {
	m = &MemoryStorage{}
	m.employees = map[string]types.Employee{}
	for _, e := range employees {
		m.employees[e.Name] = e
	}

	m.gifts = map[string]types.Gift{}
	for _, g := range gifts {
		m.gifts[g.Name] = g
	}

	return
}

func (m *MemoryStorage) AssignEmployeeToGift(employeeName, giftName string) (err error) {
	e, err := m.GetEmployeeByName(employeeName)
	if err != nil {
		return
	}
	if e.AssignedToGift != "" {
		return errors.Errorf("employee: %v is already assigned to gift: %v", employeeName, e.AssignedToGift)
	}

	g, err := m.GetGiftByName(giftName)
	if err != nil {
		return
	}
	if e.AssignedToGift != "" {
		return errors.Errorf("gift: %v is already assigned to employee: %v", employeeName, e.AssignedToGift)
	}

	e.AssignedToGift = giftName
	m.employees[employeeName] = e

	g.AssignedToEmployee = employeeName
	m.gifts[giftName] = g

	return
}

func (m *MemoryStorage) GetEmployees() (es []types.Employee, err error) {
	for _, e := range m.employees {
		es = append(es, e)
	}

	return
}

func (m *MemoryStorage) GetUnassignedEmployees() (es []types.Employee, err error) {
	for _, e := range m.employees {
		if e.AssignedToGift != "" {
			continue
		}
		es = append(es, e)
	}

	return
}

func (m *MemoryStorage) GetEmployeeByName(name string) (e types.Employee, err error) {
	e, ok := m.employees[name]
	if !ok {
		err = errors.Errorf("employee not found name: %v", name)
	}

	return
}

func (m *MemoryStorage) GetGifts() (gs []types.Gift, err error) {
	for _, g := range m.gifts {
		gs = append(gs, g)
	}

	return
}

func (m *MemoryStorage) GetUnassignedGifts() (gs []types.Gift, err error) {
	for _, g := range m.gifts {
		if g.AssignedToEmployee != "" {
			continue
		}
		gs = append(gs, g)
	}

	return
}

func (m *MemoryStorage) GetGiftByName(name string) (g types.Gift, err error) {
	g, ok := m.gifts[name]
	if !ok {
		err = errors.Errorf("gift not found name: %v", name)
	}

	return
}
