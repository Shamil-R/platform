package schema


type OrderDirective struct {
	*Directive
	argType string
}

func (d *OrderDirective) ArgType() string {
	return directiveArgument(&d.argType, d, "type")
}
