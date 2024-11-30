package lint

import (
	"errors"
	"iter"

	"github.com/goccy/go-yaml/token"
)

var commentRequireStartingSpace = errors.New("comment must start with a space")

type Comments struct {
	RequireStartingSpace bool
	IgnoreShebangs       bool
}

func (c Comments) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		if c.RequireStartingSpace {
			c.checkStartingSpace(ctx, yield)
		}
	}
}

func (c Comments) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}

func (c Comments) checkStartingSpace(ctx tokenContext, yield func(Problem) bool) error {
	if ctx.currentToken.Type != token.CommentType {
		return nil
	}

	if c.IgnoreShebangs && ctx.currentToken.Position.Line == 1 && ctx.currentToken.Position.Column == 1 {
		return nil
	}

	if len(ctx.currentToken.Value) < 2 {
		return nil
	}

	if ctx.currentToken.Value[0] != ' ' {
		problem := Problem{
			Line:   ctx.currentToken.Position.Line,
			Column: ctx.currentToken.Position.Column,
			Error:  newLintError(commentRequireStartingSpace),
		}
		if !yield(problem) {
			return nil
		}
	}

	return nil
}
