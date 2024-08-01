package score

type Filter interface {
	Check(handler string) bool
}

// ----------------------------------------------------------------------------

type includeFilter map[string]bool

func NewIncludeFilter(handlers ...string) Filter {
	flt := make(includeFilter)
	for _, handler := range handlers {
		flt[handler] = true
	}
	return flt
}

func (flt includeFilter) Check(handler string) bool {
	return flt[handler]
}

// ----------------------------------------------------------------------------

type excludeFilter map[string]bool

func NewExcludeFilter(handlers ...string) Filter {
	flt := make(excludeFilter)
	for _, handler := range handlers {
		flt[handler] = true
	}
	return flt
}

func (flt excludeFilter) Check(handler string) bool {
	return !flt[handler]
}
