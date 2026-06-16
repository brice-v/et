package action

import "fmt"

type Action int

const (
	Quit Action = iota
)

var names = map[string]Action{
	"quit": Quit,
}

func Parse(s string) (Action, error) {
	if a, ok := names[s]; ok {
		return a, nil
	}
	return -1, fmt.Errorf("unknown action %q", s)
}
