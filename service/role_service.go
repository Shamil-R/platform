package service

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/ast"
)

type RoleService interface {
	CheckRole(ctx *graphql.ResolverContext) (bool, error)
}

type roleService struct {
}

func NewRoleService() *roleService {
	return &roleService{}
}

func (s *roleService) CheckRole(ctx *graphql.ResolverContext) (bool, error) {
	fmt.Println(ctx.Object)
	var sels []string

	fieldSelections := ctx.Field.Selections

	for _, sel := range fieldSelections {
		switch sel := sel.(type) {
		case *ast.Field:
			sels = append(sels, fmt.Sprintf("%s as %s in %s", sel.Name, sel.Alias, sel.ObjectDefinition.Name))
		case *ast.InlineFragment:
			sels = append(sels, fmt.Sprintf("inline fragment on %s", sel.TypeCondition))
		case *ast.FragmentSpread:
			// fragment := reqCtx.Doc.Fragments.ForName(sel.Name)
			// sels = append(sels, fmt.Sprintf("named fragment %s on %s", sel.Name, fragment.TypeCondition))
		}
	}
	fmt.Println(sels)
	return true, nil
}

func (s *roleService) checkQuery() {}
