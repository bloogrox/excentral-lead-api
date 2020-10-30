package sender

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/sender"
)

type repo struct {
	db *gorm.DB
}

// New ...
func New(db *gorm.DB) sender.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Get(ID uint) (*sender.Sender, error) {
	var s sender.Sender

	result := r.db.First(&s, ID)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "Adapter.Sender.Get")
	}

	return &s, nil
}
