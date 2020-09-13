package postback

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"

	"gitlab.com/cpanova/excentral/domain/lead"
	"gitlab.com/cpanova/excentral/domain/postback"

	event "gitlab.com/cpanova/excentral/worker/conversion"
)

// StructHandler ...
type StructHandler struct {
	leadRepo     lead.Repo
	postbackRepo postback.Repo
}

// NewHandler ...
func NewHandler(
	leadRepo lead.Repo,
	postbackRepo postback.Repo,
) StructHandler {
	return StructHandler{
		leadRepo,
		postbackRepo,
	}
}

// Handler ...
func (s StructHandler) Handler(msg *message.Message) error {
	var e event.Event
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		log.Printf("Worker.Conversion.Postback.Handler: %s\n", err.Error())
		return nil
	}

	l, err := s.leadRepo.Get(e.LeadID)
	if err != nil {
		log.Printf("Worker.Conversion.Postback.Handler: %s\n", err.Error())
		return nil
	}

	pbs, err := s.postbackRepo.ByPID(l.PID)
	if err != nil {
		log.Printf("Worker.Conversion.Postback.Handler: %s\n", err.Error())
		return nil
	}
	for _, p := range pbs {
		url := p.URL
		url = strings.Replace(url, "{sub1}", l.Sub1, 1)
		url = strings.Replace(url, "{status}", e.Status, 1)

		_, err = http.Get(url)
		if err != nil {
			log.Printf("Worker.Conversion.Postback.Handler: %s\n", err.Error())
			return nil
		}
	}

	return nil
}
