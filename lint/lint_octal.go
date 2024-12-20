package lint

import (
	"errors"
	"iter"

	"github.com/goccy/go-yaml/token"
)

var (
	ErrOctalImplicit = errors.New("implicit octal literals are forbidden")
	ErrExplicitOctal = errors.New("explicit octal literals are forbidden")
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
		problem := problem(
			ctx.currentToken.Position.Line,
			ctx.currentToken.Position.Column,
			ErrOctalImplicit,
		)
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
		problem := problem(
			ctx.currentToken.Position.Line,
			ctx.currentToken.Position.Column,
			ErrExplicitOctal,
		)
		if !yield(problem) {
			return nil
		}
	}

	return nil
}
