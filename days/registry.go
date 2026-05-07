package days

var registry = map[int]func() Solution{}

// Register adds a day solver constructor to the package registry so main can
// instantiate solvers by day number at runtime.
func Register(day int, constructor func() Solution) {
	registry[day] = constructor
}

// Get looks up day in the registry and returns a fresh solver plus true, or nil
// and false when no solver has been registered for that day.
func Get(day int) (Solution, bool) {
	constructor, ok := registry[day]
	if !ok {
		return nil, false
	}
	return constructor(), true
}
