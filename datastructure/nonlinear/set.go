package nonlinear

type StringSet map[string]struct{}

func (set StringSet) Add(str string) {
	set[str] = struct{}{}
}

func (set StringSet) Contains(path string) bool {
	_, ok := set[path]

	return ok
}

func NewStringSet() StringSet {
	return StringSet{}
}
