package entities

type AppRoute string

const (
	GetAllCustomers         AppRoute = "getallcustomers"
	GetAllCustomersByStatus AppRoute = "getallcustomersbystatus"
	GetCustomerById         AppRoute = "getcustomerbyid"
	NewAccount              AppRoute = "newaccount"
	NewTransaction          AppRoute = "newtransaction"
)

type Role string

const (
	CLIENT Role = "user"
	ADMIN  Role = "admin"
)
