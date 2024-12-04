package lint

import (
	"errors"
	"iter"
	"strings"

	"github.com/goccy/go-yaml/token"
)

var (
	ErrBracesForbidden          = errors.New("braces are forbidden")
	ErrBracesNonEmptyForbidden  = errors.New("non empty braces are forbidden")
	ErrBracesTooFewSpaces       = errors.New("too few spaces inside braces")
	ErrBracesTooManySpaces      = errors.New("too many spaces inside braces")
	ErrBracesTooFewSpacesEmpty  = errors.New("too few spaces inside empty braces")
	ErrBracesTooManySpacesEmpty = errors.New("too many spaces inside empty braces")
)

type ForbidBraces int

const (
	ForbidBracesNone ForbidBraces = iota
	ForbidBracesAll
	ForbidBracesNonEmpty
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
		if b.Forbid == ForbidBracesAll || b.Forbid == ForbidBracesNonEmpty {
			if !b.checkBraces(ctx, yield) {
				return
			}
		}

		if !b.checkSpacesStart(ctx, yield) {
			return
		}
		if !b.checkSpacesEnd(ctx, yield) {
			return
		}
	}
}

func (b Braces) checkSpacesStart(ctx tokenContext, yield func(Problem) bool) bool {
	if ctx.nextToken == nil {
		return true
	}

	if ctx.currentToken.Type == token.MappingStartType {
		if ctx.nextToken.Type == token.MappingEndType {
			spaces := strings.Count(ctx.nextToken.Origin, " ")

			if spaces < b.MinSpacesInsideEmpty {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					ErrBracesTooFewSpacesEmpty,
				)
				if !yield(problem) {
					return false
				}
			}

			if spaces > b.MaxSpacesInsideEmpty {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					ErrBracesTooManySpacesEmpty,
				)
				if !yield(problem) {
					return false
				}
			}

			return true
		}

		leadingSpaces := leadingSpaces(ctx.nextToken.Origin)

		if leadingSpaces < b.MinSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracesTooFewSpaces,
			)
			if !yield(problem) {
				return false
			}
		}

		if leadingSpaces > b.MaxSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracesTooManySpaces,
			)
			if !yield(problem) {
				return false
			}
		}
	}

	return false
}

func (b Braces) checkSpacesEnd(ctx tokenContext, yield func(Problem) bool) bool {
	if ctx.lastToken == nil {
		return true
	}

	if ctx.currentToken.Type == token.MappingEndType {
		if ctx.lastToken.Type == token.MappingStartType {
			leadingSpaces := strings.Count(ctx.currentToken.Origin, " ")

			if leadingSpaces < b.MinSpacesInsideEmpty {
				problem := problem(
					ctx.lastToken.Position.Line,
					ctx.lastToken.Position.Column,
					ErrBracesTooFewSpacesEmpty,
				)
				if !yield(problem) {
					return false
				}
			}

			if leadingSpaces > b.MaxSpacesInsideEmpty {
				problem := problem(
					ctx.lastToken.Position.Line,
					ctx.lastToken.Position.Column,
					ErrBracesTooManySpacesEmpty,
				)
				if !yield(problem) {
					return false
				}
			}

			return true
		}

		trailingSpaces := trailingSpaces(ctx.lastToken.Origin)

		if trailingSpaces < b.MinSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracesTooFewSpaces,
			)
			if !yield(problem) {
				return false
			}
		}

		if trailingSpaces > b.MaxSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracesTooManySpaces,
			)
			if !yield(problem) {
				return false
			}
		}
	}

	return true
}

func (b Braces) checkBraces(ctx tokenContext, yield func(Problem) bool) bool {
	switch ctx.currentToken.Type {
	case token.MappingStartType:
		if ctx.currentToken.Value != "{" {
			return true
		}
	case token.MappingEndType:
		if ctx.currentToken.Value != "}" {
			return true
		}
	default:
		return true
	}

	if b.Forbid == ForbidBracesAll {
		problem := problem(
			ctx.currentToken.Position.Line,
			ctx.currentToken.Position.Column,
			ErrBracesForbidden,
		)
		if !yield(problem) {
			return false
		}

		return false
	}

	if b.Forbid == ForbidBracesNonEmpty {
		if ctx.nextToken != nil &&
			ctx.nextToken.Type == token.MappingEndType &&
			ctx.currentToken.Type == token.MappingStartType {
			spaces := strings.Count(ctx.nextToken.Origin, " ")

			if spaces > 0 {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					ErrBracesNonEmptyForbidden,
				)
				if !yield(problem) {
					return false
				}

				return false
			}
		} else if ctx.lastToken != nil &&
			ctx.lastToken.Type == token.MappingStartType &&
			ctx.currentToken.Type == token.MappingEndType {
			spaces := strings.Count(ctx.currentToken.Origin, " ")

			if spaces > 0 {
				problem := problem(
					ctx.currentToken.Position.Line,
					ctx.currentToken.Position.Column,
					ErrBracesNonEmptyForbidden,
				)
				if !yield(problem) {
					return false
				}

				return false
			}
		} else {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracesNonEmptyForbidden,
			)
			if !yield(problem) {
				return false
			}

			return false
		}
	}

	return true
}

func (b Braces) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}
