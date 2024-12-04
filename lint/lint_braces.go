package lint

import (
	"errors"
	"iter"
	"strings"

	"github.com/goccy/go-yaml/token"
)

type ForbidBraces int

const (
	ForbidBracesNone ForbidBraces = iota
	ForbidBracesAll
	ForbidBracesEmpty
)

type Braces struct {
	Forbid               ForbidBraces
	MinSpacesInside      int
	MaxSpacesInside      int
	MinSpacesInsideEmpty int
	MaxSpacesInsideEmpty int
}

func (b Braces) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		b.checkSpacesStart(ctx, yield)
		b.checkSpacesEnd(ctx, yield)
	}
}

func (b Braces) checkSpacesStart(ctx tokenContext, yield func(Problem) bool) {
	if ctx.nextToken == nil {
		return
	}

	if ctx.currentToken.Type == token.MappingStartType {
		if ctx.nextToken.Type == token.MappingEndType {
			spaces := strings.Count(ctx.nextToken.Origin, " ")

			if spaces < b.MinSpacesInsideEmpty {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					errors.New("too few spaces inside braces"),
				)
				if !yield(problem) {
					return
				}
			}

			if spaces > b.MaxSpacesInsideEmpty {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					errors.New("too many spaces inside braces"),
				)
				if !yield(problem) {
					return
				}
			}

			return
		}

		leadingSpaces := strings.LastIndexFunc(ctx.nextToken.Origin, func(r rune) bool {
			return r != ' '
		})

		if leadingSpaces < b.MinSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				errors.New("too few spaces inside braces"),
			)
			if !yield(problem) {
				return
			}
		}

		if leadingSpaces > b.MaxSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				errors.New("too many spaces inside braces"),
			)
			if !yield(problem) {
				return
			}
		}

		return
	}
}

func (b Braces) checkSpacesEnd(ctx tokenContext, yield func(Problem) bool) {
	if ctx.lastToken == nil {
		return
	}

	if ctx.currentToken.Type == token.MappingEndType {
		if ctx.lastToken.Type == token.MappingStartType {
			spaces := strings.Count(ctx.currentToken.Origin, " ")

			if spaces < b.MinSpacesInsideEmpty {
				problem := problem(
					ctx.lastToken.Position.Line,
					ctx.lastToken.Position.Column,
					errors.New("too few spaces inside braces"),
				)
				if !yield(problem) {
					return
				}
			}

			if spaces > b.MaxSpacesInsideEmpty {
				problem := problem(
					ctx.lastToken.Position.Line,
					ctx.lastToken.Position.Column,
					errors.New("too many spaces inside braces"),
				)
				if !yield(problem) {
					return
				}
			}

			return
		}

		trailingSpaces := trailingSpaces(ctx.lastToken.Origin)

		if trailingSpaces < b.MinSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				errors.New("too few spaces inside braces"),
			)
			if !yield(problem) {
				return
			}
		}

		if trailingSpaces > b.MaxSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				errors.New("too many spaces inside braces"),
			)
			if !yield(problem) {
				return
			}
		}

		return
	}
}

func (b Braces) checkBraces(ctx tokenContext, yield func(Problem) bool) {
	switch ctx.currentToken.Type {
	case token.MappingStartType:
		if ctx.currentToken.Value != "{" {
			return
		}
	case token.MappingEndType:
		if ctx.currentToken.Value != "}" {
			return
		}
	default:
		return
	}

	if b.Forbid == ForbidBracesAll {
		problem := problem(
			ctx.currentToken.Position.Line,
			ctx.currentToken.Position.Column,
			errors.New("braces are forbidden"),
		)
		if !yield(problem) {
			return
		}
	}

	if b.Forbid == ForbidBracesEmpty {
		if ctx.nextToken == nil {
			return
		}

		if ctx.nextToken.Type == token.MappingEndType {
			problem := problem(
				ctx.nextToken.Position.Line,
				ctx.nextToken.Position.Column,
				errors.New("empty braces are forbidden"),
			)
			if !yield(problem) {
				return
			}
		}
	}
}

func (b Braces) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}
