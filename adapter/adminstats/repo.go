package adminstats

import (
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/adminstats"
	"gitlab.com/cpanova/excentral/domain/lead"
)

type repo struct {
	db *gorm.DB
}

// New ...
func New(db *gorm.DB) adminstats.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) ByDay() ([]adminstats.DailyReport, error) {
	var rows []adminstats.DailyReport

	r.db.Model(&lead.Lead{}).Select("date(created_at) as day, count(*) as count").Group("date(created_at)").Order("day desc").Scan(&rows)
	// if result.Error != nil {
	// 	return nil, errors.Wrap(result.Error, "Adapter.AdminStats.ByDay")
	// }

	return rows, nil
}

func (r *repo) ByPID() ([]adminstats.PIDReport, error) {
	var rows []adminstats.PIDReport

	r.db.Model(&lead.Lead{}).Select("p_id as pid, count(*) as count").Group("p_id").Scan(&rows)
	// if result.Error != nil {
	// 	return nil, errors.Wrap(result.Error, "Adapter.AdminStats.ByPID")
	// }

	return rows, nil
}
