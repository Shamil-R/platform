package schema

const (
	ConditionEq  = "eq"
	ConditionNot = "not"
	ConditionLt  = "lt"
	ConditionLte = "lte"
	ConditionGt  = "gt"
	ConditionGte = "gte"
)

type ConditionDirective struct {
	*Directive
	argType string
}

func (d *ConditionDirective) ArgType() string {
	return directiveArgument(&d.argType, d, "type")
}
