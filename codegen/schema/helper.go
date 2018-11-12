package schema

type DirectiveName string

type ArgumentName string

func GetDirectiveArgument(d Directives, dn DirectiveName, an ArgumentName) *Directive {
	return d.Directives().ByName("name")
}
