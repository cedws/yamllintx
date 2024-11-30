package lint

import (
	"errors"
	"iter"

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

func (a Anchors) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		if ctx.currentToken.Type == token.AnchorType && ctx.nextToken != nil {
			anchorName := ctx.nextToken.Value

			if _, ok := a.declaredAnchors[anchorName]; ok && a.ForbidDuplicatedAnchors {
				problem := Problem{
					Line:   ctx.currentToken.Position.Line,
					Column: ctx.currentToken.Position.Column,
					Error:  newLintError(errors.New("anchor is duplicated")),
				}
				if !yield(problem) {
					return
				}
			}

			a.declaredAnchors[anchorName] = struct{}{}
		}

		if ctx.currentToken.Type == token.AliasType && ctx.nextToken != nil {
			anchorName := ctx.nextToken.Value

			if _, ok := a.declaredAnchors[anchorName]; !ok && a.ForbidUndeclaredAliases {
				problem := Problem{
					Line:   ctx.currentToken.Position.Line,
					Column: ctx.currentToken.Position.Column,
					Error:  newLintError(errors.New("alias references an undeclared anchor")),
				}
				if !yield(problem) {
					return
				}
			}

			a.usedAnchors[anchorName] = struct{}{}
		}

		// Final token
		if ctx.nextToken == nil && a.ForbidUnusedAnchors {
			for anchorName := range a.declaredAnchors {
				if _, ok := a.usedAnchors[anchorName]; !ok {
					problem := Problem{
						Line:   ctx.currentToken.Position.Line,
						Column: ctx.currentToken.Position.Column,
						Error:  newLintError(errors.New("anchor is declared but not used")),
					}
					if !yield(problem) {
						return
					}
				}
			}
		}
	}
}

func (a Anchors) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}