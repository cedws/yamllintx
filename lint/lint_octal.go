package lint

import (
	"errors"
	"iter"

	"github.com/goccy/go-yaml/token"
)

type Octal struct {
	ForbidImplicitOctal bool
	ForbidExplicitOctal bool
}

func (o Octal) CheckToken(ctx tokenContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {
		if o.ForbidImplicitOctal {
			o.checkImplicitOctal(ctx, yield)
		}

		if o.ForbidExplicitOctal {
			o.checkExplicitOctal(ctx, yield)
		}
	}
}

func (o Octal) CheckLine(ctx lineContext) iter.Seq[Problem] {
	return func(yield func(Problem) bool) {}
}

func (o Octal) checkImplicitOctal(ctx tokenContext, yield func(Problem) bool) {
	if ctx.currentToken.Type != token.OctetIntegerType {
		return
	}

	if len(ctx.currentToken.Value) < 2 {
		return
	}

	if ctx.currentToken.Value[0] == '0' && ctx.currentToken.Value[1] != 'o' {
		problem := Problem{
			Line:   ctx.currentToken.Position.Line,
			Column: ctx.currentToken.Position.Column,
			Error:  newLintError(errors.New("implicit octal literals are forbidden")),
		}
		if !yield(problem) {
			return
		}
	}
}

func (o Octal) checkExplicitOctal(ctx tokenContext, yield func(Problem) bool) error {
	if ctx.currentToken.Type != token.OctetIntegerType {
		return nil
	}

	if len(ctx.currentToken.Value) < 2 {
		return nil
	}

	if ctx.currentToken.Value[0] == '0' && ctx.currentToken.Value[1] == 'o' {
		problem := Problem{
			Line:   ctx.currentToken.Position.Line,
			Column: ctx.currentToken.Position.Column,
			Error:  newLintError(errors.New("explicit octal literals are forbidden")),
		}
		if !yield(problem) {
			return nil
		}
	}

	return nil
}
