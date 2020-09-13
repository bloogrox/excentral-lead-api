package lead

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/lead"
)

type repo struct {
	db *gorm.DB
	// timeout  time.Duration
}

// New ...
func New(db *gorm.DB) lead.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Get(ID uint) (*lead.Lead, error) {
	var l lead.Lead

	result := r.db.First(&l, ID)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "Adapter.Lead.Get")
	}

	return &l, nil
}

func (r *repo) Insert(l *lead.Lead) (uint, error) {
	result := r.db.Create(&l)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "Adapter.Lead.Insert")
	}

	return l.ID, nil
}
