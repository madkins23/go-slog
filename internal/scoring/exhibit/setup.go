package exhibit

import "fmt"

func Setup() error {
	if err := setupList(); err != nil {
		return fmt.Errorf("exhibit.setupList: %w", err)
	}
	if err := setupTable(); err != nil {
		return fmt.Errorf("exhibit.setupTable: %w", err)
	}
	return nil
}
