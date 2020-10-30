package partnerstats

// Repo ...
type Repo interface {
	ByDay(PID int) ([]DailyReport, error)
}
