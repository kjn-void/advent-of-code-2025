package days

var registry = map[int]func() Solution{}

func Register(day int, constructor func() Solution) {
	registry[day] = constructor
}

func Get(day int) (Solution, bool) {
	constructor, ok := registry[day]
	if !ok {
		return nil, false
	}
	return constructor(), true
}
