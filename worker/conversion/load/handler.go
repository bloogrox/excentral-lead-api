package load

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"

	advapi "gitlab.com/cpanova/excentral/ext/excentral"
	event "gitlab.com/cpanova/excentral/worker/conversion"
)

// StructHandler ...
type StructHandler struct {
	api *advapi.API
}

// NewHandler ...
func NewHandler(api *advapi.API) StructHandler {
	return StructHandler{api}
}

// Handler ...
func (s StructHandler) Handler(msg *message.Message) ([]*message.Message, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Printf("Worker.Conversion.Load: Failed load location: %s \n", err.Error())
		return nil, nil
	}
	dateStr := time.Now().In(loc).Format("2006-01-02")
	req := advapi.RequestGetConversions{
		DateFrom: dateStr + "T23:59:59+03:00",
		DateTo:   dateStr + "T00:00:00+03:00",
		Limit:    1000,
		Offset:   0,
	}
	resp, err := s.api.GetConversions(req)
	if err != nil {
		log.Printf("Worker.Conversion.Load: %s \n", err.Error())
		return nil, nil
	}

	var messages message.Messages

	for _, lead := range resp.Data.Leads {
		leadID64, err := strconv.ParseUint(lead.LeadID, 10, 32)
		if err != nil {
			log.Printf("Worker.Conversion.Load: %s\n", err.Error())
			return nil, nil
		}
		leadID := uint(leadID64)
		e := event.Event{
			LeadID:        leadID,
			DateConverted: lead.DateConverted,
			Status:        lead.Status,
		}
		b, err := json.Marshal(e)
		if err != nil {
			log.Printf("Worker.Conversion.Load: %s\n", err.Error())
			return nil, nil
		}
		msg = message.NewMessage(watermill.NewUUID(), message.Payload(b))
		messages = append(messages, msg)
	}

	return messages, nil
}
