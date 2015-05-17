package player

type Pager struct {
	From  int
	Limit int
	Total int
}

type Result struct {
	Method  string
	Url     string
	Success bool
}

type PagedResult struct {
	Pager   Pager
	Method  string
	Url     string
	Success bool
}

type Iterator struct {
	Errors chan error
	Quit   chan bool
}

func newIterator() Iterator {
	return Iterator{
		Errors: make(chan error, 1),
		Quit:   make(chan bool, 1),
	}
}
