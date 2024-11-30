package lint

import (
	"errors"
	"fmt"
	"iter"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

var lintError = errors.New("lint error")

func newLintError(err error, pos *token.Position) error {
	return fmt.Errorf("%w: %w (%d:%d)", lintError, err, pos.Line, pos.Column)
}

type sourceContext struct {
	lastToken    *token.Token
	currentToken *token.Token
	nextToken    *token.Token
}

type Linter interface {
	Check(sourceContext) error
}

type Chain []Linter

// LintAll performs linting on the entire source code and returns an iterator of all errors found.
func LintAll(src string, linters ...Linter) iter.Seq[error] {
	seqFunc := func(yield func(error) bool) {
		tokens := lexer.Tokenize(src)

		for _, lint := range linters {
			for i := 0; i < len(tokens); i++ {
				srcContext := sourceContext{
					currentToken: tokens[i],
				}

				if i >= 1 {
					srcContext.lastToken = tokens[i-1]
				}

				if i < len(tokens)-1 {
					srcContext.nextToken = tokens[i+1]
				}

				if err := lint.Check(srcContext); err != nil {
					if !yield(err) {
						return
					}
				}
			}
		}
	}

	return seqFunc
}

// Lint performs linting on the entire source code and returns the first error found.
func Lint(src string, linters ...Linter) error {
	for err := range LintAll(src, linters...) {
		return err
	}

	return nil
}
