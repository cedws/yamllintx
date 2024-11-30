package lint

import (
	"errors"

	"github.com/goccy/go-yaml/token"
)

var commentRequireStartingSpace = errors.New("comment must start with a space")

type Comments struct {
	RequireStartingSpace bool
	IgnoreShebangs       bool
}

func (c Comments) CheckToken(ctx tokenConext) error {
	if c.RequireStartingSpace {
		if err := c.checkStartingSpace(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (c Comments) CheckLine(ctx lineContext) error {
	return nil
}

func (c Comments) checkStartingSpace(ctx tokenConext) error {
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
		return newLintErrorForPosition(commentRequireStartingSpace, ctx.currentToken.Position)
	}

	return nil
}
