package persist

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"gitlab.com/cpanova/excentral/domain/conversion"

	event "gitlab.com/cpanova/excentral/worker/conversion"
)

// StructHandler ...
type StructHandler struct {
	conversionRepo conversion.Repo
}

// NewHandler ...
func NewHandler(conversionRepo conversion.Repo) StructHandler {
	return StructHandler{conversionRepo}
}

// Handler ...
func (s StructHandler) Handler(msg *message.Message) error {
	var e event.Event
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		log.Printf("Worker.Conversion.Persist.Handler: %s\n", err.Error())
		return nil
	}

	t, err := time.Parse(time.RFC3339Nano, e.DateConverted)
	if err != nil {
		log.Printf("Worker.Conversion.Persist.Handler: %s\n", err.Error())
		return nil
	}
	c := conversion.Conversion{
		LeadID:        e.LeadID,
		DateConverted: t,
		Status:        e.Status,
	}
	_, err = s.conversionRepo.Insert(&c)
	if err != nil {
		log.Printf("Worker.Conversion.Persist.Handler: %s\n", err.Error())
		return nil
	}

	return nil
}
