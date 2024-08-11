package score

import "log/slog"

type checkState int8

const (
	exclude checkState = -1
	ok                 = 0
	include            = 1
)

// ----------------------------------------------------------------------------

type Filter interface {
	Include(handler string) bool
	check(handler string) checkState
	defaultState() checkState
}

// ----------------------------------------------------------------------------

type Handler struct {
	name  string
	state checkState
}

func NewHandler(name string, defaultState checkState) Handler {
	return Handler{name: name, state: defaultState}
}

func (hs Handler) Include(handler string) bool {
	return hs.name == handler
}

func (hs Handler) String() string {
	return hs.name
}

func (hs Handler) check(handler string) checkState {
	if hs.name == handler {
		return hs.state
	} else {
		return ok
	}
}

func (hs Handler) defaultState() checkState {
	return hs.state
}

// ----------------------------------------------------------------------------

type filterCore struct {
	filters []Filter
}

func newFilterCore(defaultState checkState, filters ...any) filterCore {
	newFilter := filterCore{
		filters: make([]Filter, 0, len(filters)),
	}
	for _, filter := range filters {
		switch flt := filter.(type) {
		case Filter:
			newFilter.filters = append(newFilter.filters, flt)
		case string:
			newFilter.filters = append(newFilter.filters, NewHandler(flt, defaultState))
		default:
			slog.Error("Not a filter", "filter", filter)
		}
	}
	return newFilter
}

func (core filterCore) check(handler string) checkState {
	for _, filter := range core.filters {
		switch chkState := filter.check(handler); chkState {
		case exclude, include:
			return chkState
		}
	}

	return ok
}

func (core filterCore) included() bool {
	return true
}

// ----------------------------------------------------------------------------

type includeFilter struct {
	filterCore
}

// NewIncludeFilter returns a new "include" Filter object.
// Filters passed in may instantiate type Filter or be of type string.
func NewIncludeFilter(filters ...any) Filter {
	return &includeFilter{
		filterCore: newFilterCore(include, filters...),
	}
}

func (incl includeFilter) Include(handler string) bool {
	return incl.check(handler) == include
}

func (incl includeFilter) defaultState() checkState {
	return include
}

// ----------------------------------------------------------------------------

type excludeFilter struct {
	filterCore
}

// NewExcludeFilter returns a new "exclude" Filter object.
// Filters passed in may instantiate type Filter or be of type string.
func NewExcludeFilter(filters ...any) Filter {
	return &excludeFilter{
		filterCore: newFilterCore(exclude, filters...),
	}
}

func (excl excludeFilter) Include(handler string) bool {
	return excl.check(handler) != exclude
}

func (excl excludeFilter) defaultState() checkState {
	return include
}
