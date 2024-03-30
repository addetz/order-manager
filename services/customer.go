package jobs

type Customer struct {
	Name string `json:"name"`
}

type CustomerService struct {
	customers []*Customer
}

func NewCustomerService() *CustomerService{
	cs := &CustomerService{
		customers: make([]*Customer, 0),
	}
	cs.customers = openCustomersFile()
	return cs
} 

func (cs *CustomerService) ListCustomers() []*Customer {
	return cs.customers
}
