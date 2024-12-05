package lint

import (
	"errors"
	"iter"
	"strings"

	"github.com/goccy/go-yaml/token"
)

var (
	ErrBracketsForbidden          = errors.New("brackets are forbidden")
	ErrBracketsNonEmptyForbidden  = errors.New("non empty brackets are forbidden")
	ErrBracketsTooFewSpaces       = errors.New("too few spaces inside brackets")
	ErrBracketsTooManySpaces      = errors.New("too many spaces inside brackets")
	ErrBracketsTooFewSpacesEmpty  = errors.New("too few spaces inside empty brackets")
	ErrBracketsTooManySpacesEmpty = errors.New("too many spaces inside empty brackets")
)

type ForbidBrackets int

const (
	ForbidBracketsNone ForbidBrackets = iota
	ForbidBracketsAll
	ForbidBracketsNonEmpty
)

type Brackets struct {
	Forbid               ForbidBrackets
	MinSpacesInside      int
	MaxSpacesInside      int
	MinSpacesInsideEmpty int
	MaxSpacesInsideEmpty int
}

func (b Brackets) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		if b.Forbid == ForbidBracketsAll || b.Forbid == ForbidBracketsNonEmpty {
			if !b.checkBrackets(ctx, yield) {
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

func (b Brackets) checkSpacesStart(ctx tokenContext, yield func(Problem) bool) bool {
	if ctx.nextToken == nil {
		return true
	}

	if ctx.currentToken.Type == token.SequenceStartType {
		if ctx.nextToken.Type == token.SequenceEndType {
			spaces := strings.Count(ctx.nextToken.Origin, " ")

			if spaces < b.MinSpacesInsideEmpty {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					ErrBracketsTooFewSpacesEmpty,
				)
				if !yield(problem) {
					return false
				}
			}

			if spaces > b.MaxSpacesInsideEmpty {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					ErrBracketsTooManySpacesEmpty,
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
				ErrBracketsTooFewSpaces,
			)
			if !yield(problem) {
				return false
			}
		}

		if leadingSpaces > b.MaxSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracketsTooManySpaces,
			)
			if !yield(problem) {
				return false
			}
		}
	}

	return false
}

func (b Brackets) checkSpacesEnd(ctx tokenContext, yield func(Problem) bool) bool {
	if ctx.lastToken == nil {
		return true
	}

	if ctx.currentToken.Type == token.SequenceEndType {
		if ctx.lastToken.Type == token.SequenceStartType {
			leadingSpaces := strings.Count(ctx.currentToken.Origin, " ")

			if leadingSpaces < b.MinSpacesInsideEmpty {
				problem := problem(
					ctx.lastToken.Position.Line,
					ctx.lastToken.Position.Column,
					ErrBracketsTooFewSpacesEmpty,
				)
				if !yield(problem) {
					return false
				}
			}

			if leadingSpaces > b.MaxSpacesInsideEmpty {
				problem := problem(
					ctx.lastToken.Position.Line,
					ctx.lastToken.Position.Column,
					ErrBracketsTooManySpacesEmpty,
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
				ErrBracketsTooFewSpaces,
			)
			if !yield(problem) {
				return false
			}
		}

		if trailingSpaces > b.MaxSpacesInside {
			problem := problem(
				ctx.currentToken.Position.Line,
				ctx.currentToken.Position.Column,
				ErrBracketsTooManySpaces,
			)
			if !yield(problem) {
				return false
			}
		}
	}

	return true
}

func (b Brackets) checkBrackets(ctx tokenContext, yield func(Problem) bool) bool {
	switch ctx.currentToken.Type {
	case token.SequenceStartType:
		if ctx.currentToken.Value != "[" {
			return true
		}
	case token.SequenceEndType:
		if ctx.currentToken.Value != "]" {
			return true
		}
	default:
		return true
	}

	if b.Forbid == ForbidBracketsAll {
		problem := problem(
			ctx.currentToken.Position.Line,
			ctx.currentToken.Position.Column,
			ErrBracketsForbidden,
		)
		if !yield(problem) {
			return false
		}

		return false
	}

	if b.Forbid == ForbidBracketsNonEmpty {
		if ctx.nextToken != nil &&
			ctx.nextToken.Type == token.SequenceEndType &&
			ctx.currentToken.Type == token.SequenceStartType {
			spaces := strings.Count(ctx.nextToken.Origin, " ")

			if spaces > 0 {
				problem := problem(
					ctx.nextToken.Position.Line,
					ctx.nextToken.Position.Column,
					ErrBracketsNonEmptyForbidden,
				)
				if !yield(problem) {
					return false
				}

				return false
			}
		} else if ctx.lastToken != nil &&
			ctx.lastToken.Type == token.SequenceStartType &&
			ctx.currentToken.Type == token.SequenceEndType {
			spaces := strings.Count(ctx.currentToken.Origin, " ")

			if spaces > 0 {
				problem := problem(
					ctx.currentToken.Position.Line,
					ctx.currentToken.Position.Column,
					ErrBracketsNonEmptyForbidden,
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
				ErrBracketsNonEmptyForbidden,
			)
			if !yield(problem) {
				return false
			}

			return false
		}
	}

	return true
}

func (b Brackets) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}
