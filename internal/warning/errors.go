package warning

import "fmt"

var _ error = &Warning{}

// Error implements the error interface so all Warning objects are error type.
func (w *Warning) Error() string {
	return fmt.Sprintf("%s [%s] %s", w.Level.String(), w.Name, w.Summary)
}

// ErrorExtra returns an error object for the warning along with an extra string.
func (w *Warning) ErrorExtra(extra string) error {
	return &warningError{
		Warning: w,
		extra:   extra,
	}
}

var _ error = &warningError{}

type warningError struct {
	*Warning
	extra string
}

// Error implements the error interface for warningError objects.
// The resulting string combines the Error() result for Warning with the extra string.
func (we *warningError) Error() string {
	return fmt.Sprintf("%s: %s", we.Warning.Error(), we.extra)
}
