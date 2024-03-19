package replace

// -----------------------------------------------------------------------------

// GroupCheck returns true if the specified stack of groups (most likely the end value)
// is acceptable for processing by a ReplaceAttr function (infra.AttrFn).
type GroupCheck func(groups []string) bool

// -----------------------------------------------------------------------------

var _ GroupCheck = TopCheck

// TopCheck is a GroupCheck that returns true if there are no groups in the stack,
// (meaning the attribute being evaluated by a ReplaceAttr function
func TopCheck(groups []string) bool {
	return len(groups) == 0
}

// -----------------------------------------------------------------------------

// Current returns a GroupCheck that will return true if the top group on the stack
// (the highest indexed group in the groups array) matches the specified name.
func Current(name string) GroupCheck {
	return func(groups []string) bool {
		return len(groups) > 0 && groups[len(groups)-1] == name
	}
}
