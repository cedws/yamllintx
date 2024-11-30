package lint

import (
	"errors"

	"github.com/goccy/go-yaml/token"
)

type Octal struct {
	ForbidImplicitOctal bool
	ForbidExplicitOctal bool
}

func (o Octal) CheckToken(ctx tokenConext) error {
	if o.ForbidImplicitOctal {
		if err := o.checkImplicitOctal(ctx); err != nil {
			return err
		}
	}

	if o.ForbidExplicitOctal {
		if err := o.checkExplicitOctal(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (o Octal) CheckLine(ctx lineContext) error {
	return nil
}

func (o Octal) checkImplicitOctal(ctx tokenConext) error {
	if ctx.currentToken.Type != token.OctetIntegerType {
		return nil
	}

	if len(ctx.currentToken.Value) < 2 {
		return nil
	}

	if ctx.currentToken.Value[0] == '0' && ctx.currentToken.Value[1] != 'o' {
		return newLintErrorForPosition(errors.New("implicit octal literals are forbidden"), ctx.currentToken.Position)
	}

	return nil
}

func (o Octal) checkExplicitOctal(ctx tokenConext) error {
	if ctx.currentToken.Type != token.OctetIntegerType {
		return nil
	}

	if len(ctx.currentToken.Value) < 2 {
		return nil
	}

	if ctx.currentToken.Value[0] == '0' && ctx.currentToken.Value[1] == 'o' {
		return newLintErrorForPosition(errors.New("explicit octal literals are forbidden"), ctx.currentToken.Position)
	}

	return nil
}
