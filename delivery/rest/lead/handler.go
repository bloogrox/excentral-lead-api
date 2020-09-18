package lead

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/cpanova/excentral/ext/excentral"

	"gitlab.com/cpanova/excentral/domain/advertiser"
	"gitlab.com/cpanova/excentral/domain/lead"
)

const pid = 1

// Handler ...
type Handler interface {
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	adv      advertiser.Service
	leadRepo lead.Repo
}

// NewHandler ...
func NewHandler(
	adv advertiser.Service,
	leadRepo lead.Repo,
) Handler {
	return &handler{
		adv:      adv,
		leadRepo: leadRepo,
	}
}

// Request ...
type Request struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Language  string `json:"language"`
	Country   string `json:"country"`
	Source    string `json:"source"`
	Sub1      string `json:"sub1"`
	PID       int    `json:"pid"`
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	var data Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	advReq := advertiser.Request{
		Email:     data.Email,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Phone:     data.Phone,
		Language:  data.Language,
		Country:   data.Country,
		Source:    data.Source,
	}
	advLead, err := h.adv.CreateLead(advReq)

	if err != nil {
		switch errors.Cause(err) {
		case
			excentral.ErrRequiredField,
			excentral.ErrFieldNotValid,
			excentral.ErrInvalidData,
			excentral.ErrDuplicateEmail,
			excentral.ErrRestrictedCountry:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		case
			excentral.ErrAPIAccess,
			excentral.ErrVerificationParametersMissing,
			excentral.ErrWrongChecksum,
			excentral.ErrUnknown:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	l := lead.Lead{
		RemoteID:  advLead.ID,
		Email:     data.Email,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Phone:     data.Phone,
		Language:  data.Language,
		Country:   data.Country,
		Source:    data.Source,
		Sub1:      data.Sub1,
		PID:       pid,
	}
	_, err = h.leadRepo.Insert(&l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(""))
}
