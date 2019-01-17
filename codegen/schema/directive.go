package schema

import (
	"github.com/vektah/gqlparser/ast"
)

const (
	DirectivePrimary   = "primary"
	DirectiveUnique    = "unique"
	DirectiveIdentity  = "identity"
	DirectiveValidate  = "validate"
	DirectiveTable     = "table"
	DirectiveField     = "field"
	DirectiveRelation  = "relation"
	DirectiveInput     = "input"
	DirectiveCondition = "condition"
	DirectiveOrder 	   = "order"
	DirectiveTimestamp = "timestamp"
	DirectiveSoftDelete= "softDelete"

	InputDirectiveCreateOneWithout = "create_one_without"
)

type Directives interface {
	Directives() DirectiveList
}

type Directive struct {
	*ast.Directive
	arguments ArgumentList
}

func (d *Directive) IsPrimary() bool {
	return isPrimaryDirective(d)
}

func (d *Directive) IsUnique() bool {
	return isUniqueDirective(d)
}

func (d *Directive) IsIndentity() bool {
	return isIdentityDirective(d)
}

func (d *Directive) IsValidate() bool {
	return isValidateDirective(d)
}

func (d *Directive) IsTable() bool {
	return isTableDirective(d)
}

func (d *Directive) IsField() bool {
	return isFieldDirective(d)
}

func (d *Directive) IsTimestamp() bool {
	return isTimestampDirective(d)
}

func (d *Directive) IsSoftDelete() bool {
	return isSoftDeleteDirective(d)
}

func (f *Directive) Arguments() ArgumentList {
	if f.arguments != nil {
		return f.arguments
	}

	f.arguments = make(ArgumentList, len(f.Directive.Arguments))
	for i, arg := range f.Directive.Arguments {
		f.arguments[i] = &Argument{Argument: arg}
	}

	return f.arguments
}

type ValidateDirective struct {
	*Directive
	argMax string
}

func (d *ValidateDirective) ArgMax() string {
	return directiveArgument(&d.argMax, d, "max")
}

type TableDirective struct {
	*Directive
	argName string
}

func (d *TableDirective) ArgName() string {
	return directiveArgument(&d.argName, d, "name")
}

type FieldDirective struct {
	*Directive
	argName string
}

func (d *FieldDirective) ArgName() string {
	return directiveArgument(&d.argName, d, "name")
}

type RelationDirective struct {
	*Directive
	argObject     string
	argField      string
	argForeignKey string
	argTable 	  string
	argOwnerKey	  string
}

func (d *RelationDirective) ArgObject() string {
	return directiveArgument(&d.argObject, d, "object")
}

func (d *RelationDirective) ArgField() string {
	return directiveArgument(&d.argField, d, "field")
}

func (d *RelationDirective) ArgForeignKey() string {
	return directiveArgument(&d.argForeignKey, d, "foreignKey")
}

func (d *RelationDirective) ArgTable() string {
	return directiveArgument(&d.argTable, d, "table")
}

func (d *RelationDirective) ArgOwnerKey() string {
	return directiveArgument(&d.argOwnerKey, d, "ownerKey")
}

type InputDirective struct {
	*Directive
}

func (d *InputDirective) IsCreateOneWithout() bool {
	return d.Arguments().ByName("name").Value().Raw == InputDirectiveCreateOneWithout
}

type TimestampDirective struct {
	*Directive
	argDisable     string
	argCreateField string
	argUpdateField string
}

func (d *TimestampDirective) ArgDisable() string {
	return directiveArgument(&d.argDisable, d, "disable")
}

func (d *TimestampDirective) ArgCreateField() string {
	return directiveArgument(&d.argCreateField, d, "createField")
}

func (d *TimestampDirective) ArgUpdateField() string {
	return directiveArgument(&d.argUpdateField, d, "updateField")
}

type SoftDeleteDirective struct {
	*Directive
	argDisable     string
	argDeleteField string
}

func (d *SoftDeleteDirective) ArgDisable() string {
	return directiveArgument(&d.argDisable, d, "disable")
}

func (d *SoftDeleteDirective) ArgDeleteField() string {
	return directiveArgument(&d.argDeleteField, d, "deleteField")
}

type DirectiveList []*Directive

func (l DirectiveList) HasPrimary() bool {
	return hasDirective(l, isPrimaryDirective)
}

func (l DirectiveList) HasUnique() bool {
	return hasDirective(l, isUniqueDirective)
}

func (l DirectiveList) HasIndentity() bool {
	return hasDirective(l, isIdentityDirective)
}

func (l DirectiveList) HasValidate() bool {
	return hasDirective(l, isValidateDirective)
}

func (l DirectiveList) HasTable() bool {
	return hasDirective(l, isTableDirective)
}

func (l DirectiveList) HasField() bool {
	return hasDirective(l, isFieldDirective)
}

func (l DirectiveList) Primary() *Directive {
	return firstDirective(l, isPrimaryDirective)
}

func (l DirectiveList) Unique() *Directive {
	return firstDirective(l, isUniqueDirective)
}

func (l DirectiveList) Identity() *Directive {
	return firstDirective(l, isIdentityDirective)
}

func (l DirectiveList) Validate() *ValidateDirective {
	directive := firstDirective(l, isValidateDirective)
	if directive == nil {
		return nil
	}
	return &ValidateDirective{Directive: directive}
}

func (l DirectiveList) Table() *TableDirective {
	directive := firstDirective(l, isTableDirective)
	if directive == nil {
		return nil
	}
	return &TableDirective{Directive: directive}
}

func (l DirectiveList) Field() *FieldDirective {
	directive := firstDirective(l, isFieldDirective)
	if directive == nil {
		return nil
	}
	return &FieldDirective{Directive: directive}
}

func (l DirectiveList) Timestamp() *TimestampDirective {
	directive := firstDirective(l, isTimestampDirective)
	if directive == nil {
		return nil
	}
	return &TimestampDirective{Directive: directive}
}

func (l DirectiveList) SoftDelete() *SoftDeleteDirective {
	directive := firstDirective(l, isSoftDeleteDirective)
	if directive == nil {
		return nil
	}
	return &SoftDeleteDirective{Directive: directive}
}

func (l DirectiveList) Relation() *RelationDirective {
	directive := firstDirective(l, isRelationDirective)
	if directive == nil {
		return nil
	}
	return &RelationDirective{Directive: directive}
}

func (l DirectiveList) Input() *InputDirective {
	directive := firstDirective(l, isInputDirective)
	if directive == nil {
		return nil
	}
	return &InputDirective{Directive: directive}
}

func (l DirectiveList) Condition() *ConditionDirective {
	directive := firstDirective(l, isConditionDirective)
	if directive == nil {
		return nil
	}
	return &ConditionDirective{Directive: directive}
}

func (l DirectiveList) OrderBy() *OrderDirective {
	directive := firstDirective(l, isOrderDirective)
	if directive == nil {
		return nil
	}
	return &OrderDirective{Directive: directive}
}

func (l DirectiveList) ByName(name string) *Directive {
	return firstDirective(l, byNameDirective(name))
}

type directiveFilter func(directive *Directive) bool

func isPrimaryDirective(directive *Directive) bool {
	return directive.Name == DirectivePrimary
}

func isUniqueDirective(directive *Directive) bool {
	return directive.Name == DirectiveUnique
}

func isIdentityDirective(directive *Directive) bool {
	return directive.Name == DirectiveIdentity
}

func isValidateDirective(directive *Directive) bool {
	return directive.Name == DirectiveValidate
}

func isTableDirective(directive *Directive) bool {
	return directive.Name == DirectiveTable
}

func isFieldDirective(directive *Directive) bool {
	return directive.Name == DirectiveField
}

func isRelationDirective(directive *Directive) bool {
	return directive.Name == DirectiveRelation
}

func isInputDirective(directive *Directive) bool {
	return directive.Name == DirectiveInput
}

func isConditionDirective(directive *Directive) bool {
	return directive.Name == DirectiveCondition
}

func isOrderDirective(directive *Directive) bool {
	return directive.Name == DirectiveOrder
}

func isTimestampDirective(directive *Directive) bool {
	return directive.Name == DirectiveTimestamp
}

func isSoftDeleteDirective(directive *Directive) bool {
	return directive.Name == DirectiveSoftDelete
}

func byNameDirective(name string) directiveFilter {
	return func(directive *Directive) bool {
		return directive.Name == name
	}
}

func hasDirective(list DirectiveList, filter directiveFilter) bool {
	for _, directive := range list {
		if filter(directive) {
			return true
		}
	}
	return false
}

func firstDirective(list DirectiveList, filter directiveFilter) *Directive {
	for _, directive := range list {
		if filter(directive) {
			return directive
		}
	}
	return nil
}

func directiveArgument(current *string, a Arguments, name string) string {
	if len(*current) == 0 {
		arg := a.Arguments().ByName(name)
		if arg == nil {
			panic("argument does not exist")
		}
		val := arg.Value()
		if val == nil {
			panic("value does not exist")
		}
		*current = val.Raw
	}
	return *current
}
