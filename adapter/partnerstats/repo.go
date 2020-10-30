package partnerstats

import (
	"gorm.io/gorm"

	"gitlab.com/cpanova/excentral/domain/lead"
	"gitlab.com/cpanova/excentral/domain/partnerstats"
)

type repo struct {
	db *gorm.DB
}

// New ...
func New(db *gorm.DB) partnerstats.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) ByDay(PID int) ([]partnerstats.DailyReport, error) {
	var rows []partnerstats.DailyReport

	r.db.Model(&lead.Lead{}).
		Select(
			`
			TO_CHAR(created_at::date, 'yyyy-mm-dd') as day,
			count(*) as leads
			`).
		Where("p_id = ?", PID).
		Group("day").
		Order("day desc").
		Scan(&rows)

	return rows, nil
}
