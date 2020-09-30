package partner

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/partner"
)

type repo struct {
	db *gorm.DB
	// timeout  time.Duration
}

// New ...
func New(db *gorm.DB) partner.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Get(ID int) (*partner.Partner, error) {
	var p partner.Partner

	result := r.db.First(&p, ID)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "Adapter.Partner.Get")
	}

	return &p, nil
}
