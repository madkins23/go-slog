package warning

import "fmt"

var _ error = &Warning{}

func (w *Warning) Error() string {
	return fmt.Sprintf("%s [%s] %s", w.Level.String(), w.Name, w.Summary)
}

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

func (we *warningError) Error() string {
	return fmt.Sprintf("%s: %s", we.Warning.Error(), we.extra)
}
