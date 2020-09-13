package postback

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/postback"
)

type repo struct {
	db *gorm.DB
	// timeout  time.Duration
}

// New ...
func New(db *gorm.DB) postback.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) ByPID(pid int) ([]postback.Postback, error) {
	var postbacks []postback.Postback

	result := r.db.Where("pid = ?", pid).Find(&postbacks)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "Adapter.Postback.ByPID")
	}

	return postbacks, nil
}
