package schema

const (
	ACTION_UNDEFINED  = "undefined"
	ACTION_CREATE     = "create"
	ACTION_UPDATE     = "update"
	ACTION_DELETE     = "delete"
	ACTION_UPSERT     = "upsert"
	ACTION_ITEM       = "item"
	ACTION_COLLECTION = "collection"
	ACTION_RELATION   = "relation"
)

type Action struct {
	*FieldDefinition
	Action string
}

type ActionList []*Action

type actionListFilter func(field *Action) bool
