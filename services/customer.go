package jobs

import (
	"fmt"

	"github.com/google/uuid"
)

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Note string `json:"note"`
}

type CustomerService struct {
	customers map[string]*Customer
	filepath  string
}

func NewCustomerService(filepath string) *CustomerService {
	cs := &CustomerService{
		customers: make(map[string]*Customer, 0),
		filepath:  filepath,
	}
	for _, c := range openCustomersFile(filepath) {
		cs.customers[c.ID] = c
	}
	return cs
}

func (cs *CustomerService) ListCustomers() []*Customer {
	csList := make([]*Customer, 0)
	for _, c := range cs.customers {
		csList = append(csList, c)
	}
	return csList
}

func (cs *CustomerService) AddCustomer(cust *Customer) {
	id := uuid.New().String()
	cust.ID = id
	cs.customers[id] = cust
	cs.exportCustomers()
}

func (cs *CustomerService) DeleteCustomer(id string) {
	delete(cs.customers, id)
	cs.exportCustomers()
}

func (cs *CustomerService) UpdateCustomer(id string, newCust *Customer) error {
	curr, ok := cs.customers[id]
	if !ok {
		return fmt.Errorf("customer %s not found", id)
	}

	if newCust.Name != "" {
		curr.Name = newCust.Name
	}

	if newCust.Note != "" {
		curr.Name = newCust.Note
	}

	cs.customers[id] = curr
	cs.exportCustomers()
	return nil
}

func (cs *CustomerService) SearchCustomer(name string) (*Customer, error) {
	for _, c := range cs.customers {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, fmt.Errorf("no customer named %s found", name)
}

func (cs *CustomerService) exportCustomers() {
	writeCustomersFile(cs.filepath, cs.ListCustomers())
}
