package fakeadvertiser

import (
	"github.com/pkg/errors"
	"gitlab.com/cpanova/excentral/ext/excentral"

	adv "gitlab.com/cpanova/excentral/domain/advertiser"
)

type fake struct{}

// New ...
func New() adv.Service {
	return &fake{}
}

func (fake) CreateLead(adv.Request) (*adv.Lead, error) {
	return nil, errors.Wrap(excentral.ErrUnknown, "fakeadapter")
	// return &adv.Lead{
	// 	ID:     1,
	// 	Status: "Completed",
	// }, nil
}
