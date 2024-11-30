package lint

import (
	"errors"

	"github.com/goccy/go-yaml/token"
)

type Anchors struct {
	ForbidUndeclaredAliases bool
	ForbidDuplicatedAnchors bool
	ForbidUnusedAnchors     bool

	declaredAnchors map[string]struct{}
	usedAnchors     map[string]struct{}
}

func NewAnchors() Anchors {
	return Anchors{
		declaredAnchors: make(map[string]struct{}),
		usedAnchors:     make(map[string]struct{}),
	}
}

func (a Anchors) CheckToken(ctx tokenConext) error {
	if ctx.currentToken.Type == token.AnchorType && ctx.nextToken != nil {
		anchorName := ctx.nextToken.Value

		if _, ok := a.declaredAnchors[anchorName]; ok && a.ForbidDuplicatedAnchors {
			return newLintErrorForPosition(errors.New("anchor is duplicated"), ctx.currentToken.Position)
		}

		a.declaredAnchors[anchorName] = struct{}{}
	}

	if ctx.currentToken.Type == token.AliasType && ctx.nextToken != nil {
		anchorName := ctx.nextToken.Value

		if _, ok := a.declaredAnchors[anchorName]; !ok && a.ForbidUndeclaredAliases {
			return newLintErrorForPosition(errors.New("alias must reference an anchor"), ctx.currentToken.Position)
		}

		a.usedAnchors[anchorName] = struct{}{}
	}

	// Final token
	if ctx.nextToken == nil && a.ForbidUnusedAnchors {
		for anchorName := range a.declaredAnchors {
			if _, ok := a.usedAnchors[anchorName]; !ok {
				return newLintError(errors.New("anchor is declared but not used"), ctx.currentToken.Position.Line, ctx.currentToken.Position.Column)
			}
		}
	}

	return nil
}

func (a Anchors) CheckLine(ctx lineContext) error {
	return nil
}
