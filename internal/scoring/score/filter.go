package score

import "log/slog"

type Filter interface {
	Keep(handler string) bool
}

// ----------------------------------------------------------------------------

type Group struct {
	filters map[string]bool
}

// NewFilterGroup creates a new Group object.
// Grouped items must be either string or *Group.
func NewFilterGroup(items ...any) *Group {
	newFilter := &Group{
		filters: make(map[string]bool, len(items)),
	}
	for _, filter := range items {
		switch flt := filter.(type) {
		case string:
			newFilter.filters[flt] = true
		case *Group:
			for key, value := range flt.filters {
				if value {
					newFilter.filters[key] = true
				} else {
					delete(newFilter.filters, key)
				}
			}
		default:
			slog.Error("Not a filter", "filter", flt)
		}
	}
	return newFilter
}

func (g *Group) keep(handler string) bool {
	return g.filters[handler]
}

// ----------------------------------------------------------------------------

type includeFilter struct {
	*Group
}

// NewIncludeFilter returns a new "include" Filter object.
// Filtered items must be either string or *Group.
func NewIncludeFilter(items ...any) Filter {
	return &includeFilter{
		Group: NewFilterGroup(items...),
	}
}

func (incl includeFilter) Keep(handler string) bool {
	return incl.keep(handler)
}

// ----------------------------------------------------------------------------

type excludeFilter struct {
	*Group
}

// NewExcludeFilter returns a new "exclude" Filter object.
// Filtered items must be either string or *Group.
func NewExcludeFilter(items ...any) Filter {
	return &excludeFilter{
		Group: NewFilterGroup(items...),
	}
}

func (excl excludeFilter) Keep(handler string) bool {
	return !excl.keep(handler)
}
