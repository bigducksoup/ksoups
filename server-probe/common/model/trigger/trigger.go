package trigger

type Trigger struct {
	Id        string
	Name      string
	Next      []Trigger
	Fuse      bool
	Shortcuts []string
}
