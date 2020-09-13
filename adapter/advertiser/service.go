package advertiser

import (
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/cpanova/excentral/ext/excentral"

	"gitlab.com/cpanova/excentral/domain/advertiser"
)

// Service ...
type Service struct {
	api *excentral.API
}

// New ...
func New(api *excentral.API) *Service {
	return &Service{api}
}

// CreateLead ...
func (s *Service) CreateLead(r advertiser.Request) (*advertiser.Lead, error) {
	req := excentral.RequestCreateLead{
		Email:     r.Email,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Phone:     r.Phone,
		Language:  r.Language,
		Country:   r.Country,
		Marker:    r.Source,
	}

	resp, err := s.api.CreateLead(req)
	if err != nil {
		return nil, errors.Wrap(err, "Adapter.Service.CreateLead")
	}

	leadID, err := strconv.Atoi(resp.Data.LeadID)
	if err != nil {
		return nil, errors.Wrap(err, "Adapter.Service.CreateLead")
	}

	lead := advertiser.Lead{
		ID:     uint(leadID),
		Status: resp.Data.Status,
	}

	return &lead, nil
}
