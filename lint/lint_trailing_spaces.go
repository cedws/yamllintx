package lint

import (
	"errors"
	"strings"
)

type TrailingSpaces struct{}

func (t TrailingSpaces) CheckToken(ctx tokenConext) error {
	return nil
}

func (t TrailingSpaces) CheckLine(ctx lineContext) error {
	lastNonSpace := strings.LastIndexFunc(ctx.currentLine, func(r rune) bool {
		return r != ' '
	})

	if lastNonSpace != len(ctx.currentLine)-1 {
		return newLintError(errors.New("trailing spaces are forbidden"), ctx.currentLineNumber, lastNonSpace+1)
	}

	return nil
}
