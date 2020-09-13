package process

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"gorm.io/gorm"

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
func (s StructHandler) Handler(msg *message.Message) ([]*message.Message, error) {
	var e event.Event
	err := json.Unmarshal(msg.Payload, &e)
	if err != nil {
		log.Printf("Worker.Conversion.Process: %s\n", err.Error())
		return nil, nil
	}
	_, err = s.conversionRepo.Get(e.LeadID)
	if err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			b, jsonErr := json.Marshal(e)
			if jsonErr != nil {
				log.Printf("Worker.Conversion.Process: %s\n", err.Error())
				return nil, nil
			}

			return message.Messages{
					message.NewMessage(
						watermill.NewUUID(),
						message.Payload(b),
					),
				},
				nil
		}

		log.Printf("Worker.Conversion.Process: %s\n", err.Error())
		return nil, nil
	}

	return nil, nil
}
