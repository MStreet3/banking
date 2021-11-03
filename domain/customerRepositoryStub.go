package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{"2000", "Steve", "1978-12-15", "Delhi", "110075", "1"},
		{"2001", "Arian", "1988-05-21", "Newburgh, NY", "12550", "1"},
	}
	return CustomerRepositoryStub{
		customers: customers,
	}
}
