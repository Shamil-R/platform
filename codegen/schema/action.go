package schema

const (
	ACTION_UNDEFINED         = "undefined"
	ACTION_CREATE            = "create"
	ACTION_UPDATE            = "update"
	ACTION_DELETE            = "delete"
	ACTION_UPSERT            = "upsert"
	ACTION_ITEM              = "item"
	ACTION_COLLECTION        = "collection"
	ACTION_RELATION          = "relation"
	ACTION_DELETE_MANY       = "deleteMany"
	ACTION_UPDATE_MANY       = "updateMany"
	ACTION_FORCE_DELETE      = "forceDelete"
	ACTION_FORCE_DELETE_MANY = "forceDeleteMany"
	ACTION_RESTORE           = "restore"
	ACTION_RESTORE_MANY      = "restoreMany"
)

type Action struct {
	*FieldDefinition
	Action string
}

type ActionList []*Action

type actionListFilter func(field *Action) bool
