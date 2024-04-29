package replace

// -----------------------------------------------------------------------------

// GroupCheckFn returns true if the specified stack of groups (most likely the end value)
// is acceptable for processing by a ReplaceAttr function (infra.AttrFn).
// Use the function pointer, returns true if the group stack is empty.
type GroupCheckFn func(groups []string) bool

// -----------------------------------------------------------------------------

var _ GroupCheckFn = TopCheck

// TopCheck is a GroupCheckFn that returns true if there are no groups in the stack,
// (meaning the attribute being evaluated by a ReplaceAttr function
// Use the result of executing this function with the current groups stack.
func TopCheck(groups []string) bool {
	return len(groups) == 0
}

// -----------------------------------------------------------------------------

// Current returns a GroupCheckFn that will return true if the top group on the stack
// (the highest indexed group in the groups array) matches the specified name.
// Use the result of executing this function with a specified group name.
// The resulting function will return true if that name is at the top
// of the groups stack (highest indexed item in the array).
// An empty stack always returns false.
func Current(name string) GroupCheckFn {
	return func(groups []string) bool {
		return len(groups) > 0 && groups[len(groups)-1] == name
	}
}
