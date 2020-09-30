package advertiser

// Request ...
type Request struct {
	Email     string
	FirstName string
	LastName  string
	Phone     string
	Language  string
	Country   string
	Source    string
	Campaign  string
}

// Lead ...
type Lead struct {
	ID     uint
	Status string
}

// Service ...
type Service interface {
	CreateLead(Request) (*Lead, error)
}
