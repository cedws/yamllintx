package lint

import (
	"errors"
	"iter"

	"github.com/goccy/go-yaml/token"
)

type anchors struct {
	AnchorOpts
	declaredAnchors map[string]token.Token
	usedAnchors     map[string]struct{}
}

type AnchorOpts struct {
	ForbidUndeclaredAliases bool
	ForbidDuplicatedAnchors bool
	ForbidUnusedAnchors     bool
}

func Anchors(opts AnchorOpts) Linter {
	return anchors{
		AnchorOpts:      opts,
		declaredAnchors: make(map[string]token.Token),
		usedAnchors:     make(map[string]struct{}),
	}
}

func (a anchors) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		if ctx.currentToken.Type == token.AnchorType && ctx.nextToken != nil {
			anchorName := ctx.nextToken.Value

			if _, ok := a.declaredAnchors[anchorName]; ok && a.ForbidDuplicatedAnchors {
				problem := problem(
					ctx.currentToken.Position.Line,
					ctx.currentToken.Position.Column,
					errors.New("anchor is duplicated"),
				)
				if !yield(problem) {
					return
				}
			}

			a.declaredAnchors[anchorName] = *ctx.currentToken
		}

		if ctx.currentToken.Type == token.AliasType && ctx.nextToken != nil {
			anchorName := ctx.nextToken.Value

			if _, ok := a.declaredAnchors[anchorName]; !ok && a.ForbidUndeclaredAliases {
				problem := problem(
					ctx.currentToken.Position.Line,
					ctx.currentToken.Position.Column,
					errors.New("alias references an undeclared anchor"),
				)
				if !yield(problem) {
					return
				}
			}

			a.usedAnchors[anchorName] = struct{}{}
		}

		// Final token
		if a.ForbidUnusedAnchors && ctx.nextToken == nil {
			for anchorName, anchor := range a.declaredAnchors {
				if _, ok := a.usedAnchors[anchorName]; !ok {
					problem := problem(
						anchor.Position.Line,
						anchor.Position.Column,
						errors.New("anchor is declared but not used"),
					)
					if !yield(problem) {
						return
					}
				}
			}
		}
	}
}

func (a anchors) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}
