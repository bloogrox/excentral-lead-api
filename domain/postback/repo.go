package postback

// Repo ...
type Repo interface {
	ByPID(int) ([]Postback, error)
}
