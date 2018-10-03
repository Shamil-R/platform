package validate

type Data struct {
	Rule  string
	Value string
}

type Validate interface {
	Validate(d *Data) error
}

type validate struct{}

func (v validate) Validate(d *Data) error {
	return nil
}
