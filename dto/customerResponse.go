package dto

type CustomerStatus string

const (
	ACTIVE   CustomerStatus = "active"
	INACTIVE CustomerStatus = "inactive"
)

func (c CustomerStatus) StatusAsQueryParam() int {
	var s int
	switch c {
	case ACTIVE:
		s = 1
	case INACTIVE:
		s = 0
	}
	return s
}

type CustomerResponse struct {
	Id          string         `json:"customer_id"`
	Name        string         `json:"full_name"`
	City        string         `json:"city"`
	Zipcode     string         `json:"zipcode"`
	DateOfBirth string         `json:"date_of_birth"`
	Status      CustomerStatus `json:"status"`
}
