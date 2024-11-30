package lint

import (
	"errors"
	"strings"

	"github.com/goccy/go-yaml/token"
)

var hypensMaxSpacesAfter = errors.New("too many spaces after hypen")

type Hyphens struct {
	MaxSpacesAfter int
}

func (h Hyphens) CheckToken(ctx tokenConext) error {
	if h.MaxSpacesAfter > 0 {
		if err := h.checkMaxSpacesAfter(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (h Hyphens) CheckLine(ctx lineContext) error {
	return nil
}

func (h Hyphens) checkMaxSpacesAfter(ctx tokenConext) error {
	if ctx.currentToken.Type != token.SequenceEntryType || ctx.nextToken == nil {
		return nil
	}

	origin := ctx.nextToken.Origin
	firstNonSpace := strings.IndexFunc(origin, func(r rune) bool {
		return r != ' '
	})
	if firstNonSpace == -1 {
		firstNonSpace = len(origin)
	}
	if firstNonSpace > h.MaxSpacesAfter {
		return newLintErrorForPosition(hypensMaxSpacesAfter, ctx.nextToken.Position)
	}

	return nil
}
