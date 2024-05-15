package setup

import "fmt"

type Index struct {
	lookup map[string]Item
}

func (i *Index) Has(name string) bool {
	_, found := i.lookup[name]
	return found
}

func (i *Index) Lookup(name string) Item {
	return i.lookup[name]
}

func (i *Index) Register(name string, item Item) {
	i.lookup[name] = item
}

func (i *Index) Setup() error {
	for name, item := range i.lookup {
		if err := item.Setup(); err != nil {
			return fmt.Errorf("setup '%s': %w", name, err)
		}
	}
	return nil
}
