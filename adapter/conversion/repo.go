package conversion

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/conversion"
)

type repo struct {
	db *gorm.DB
	// timeout  time.Duration
}

// New ...
func New(db *gorm.DB) conversion.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Get(LeadID uint) (*conversion.Conversion, error) {
	var c conversion.Conversion

	result := r.db.First(&c, LeadID)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "Adapter.Conversion.Get")
	}

	return &c, nil
}

func (r *repo) Insert(c *conversion.Conversion) (uint, error) {
	result := r.db.Create(&c)
	if result.Error != nil {
		return 0, errors.Wrap(result.Error, "Adapter.Conversion.Insert")
	}

	return c.LeadID, nil
}
