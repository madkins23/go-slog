package setup

type Item interface {
	Name() string
	Setup() error
}

// -----------------------------------------------------------------------------

type ItemCore struct {
	name string
}

func (ic *ItemCore) Name() string {
	return ic.name
}

func (ic *ItemCore) Setup() error {
	return nil
}
