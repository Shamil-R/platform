package schema

type Action struct {
	Action          string
	FieldDefinition *FieldDefinition
	Definition      *Definition
}

func (a *Action) IsRelation() bool {
	return a.Action == ACTION_RELATION
}

type ActionList []*Action

type actionListFilter func(field *Action) bool

func (l ActionList) size() int {
	return len(l)
}

func (l ActionList) filter(filter actionListFilter) ActionList {
	actions := make(ActionList, 0, len(l))
	for _, action := range l {
		if filter(action) {
			actions = append(actions, action)
		}
	}
	return actions
}

func (l ActionList) first(filter actionListFilter) *Action {
	r := l.filter(filter)
	if r.size() == 0 {
		return nil
	}
	return r[0]
}

func (l ActionList) ByAction(action string) *Action {
	fn := func(a *Action) bool {
		return a.Action == action
	}
	return l.first(fn)
}
