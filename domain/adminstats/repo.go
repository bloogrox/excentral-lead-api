package adminstats

// Repo ...
type Repo interface {
	ByDay() ([]DailyReport, error)
	ByPID() ([]PIDReport, error)
}
