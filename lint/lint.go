package lint

import (
	"bufio"
	"errors"
	"fmt"
	"iter"
	"strings"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

var lintError = errors.New("lint error")

func newLintError(err error, line, column int) error {
	return fmt.Errorf("%w: %w (%d:%d)", lintError, err, line, column)
}

func newLintErrorForPosition(err error, pos *token.Position) error {
	return newLintError(err, pos.Line, pos.Column)
}

type tokenConext struct {
	lastToken    *token.Token
	currentToken *token.Token
	nextToken    *token.Token
}

type lineContext struct {
	currentLine       string
	currentLineNumber int
}

type Linter interface {
	CheckToken(tokenConext) error
	CheckLine(lineContext) error
}

type Chain []Linter

// LintAll performs linting on the entire source code and returns an iterator of all errors found.
func LintAll(src string, linters ...Linter) iter.Seq[error] {
	tokens := lexer.Tokenize(src)
	tokens.Dump()

	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	seqFunc := func(yield func(error) bool) {
		for _, lint := range linters {
			for i := 0; i < len(lines); i++ {
				lineContext := lineContext{
					currentLine:       lines[i],
					currentLineNumber: i + 1,
				}

				if err := lint.CheckLine(lineContext); err != nil {
					if !yield(err) {
						return
					}
				}
			}

			for i := 0; i < len(tokens); i++ {
				srcContext := tokenConext{
					currentToken: tokens[i],
				}

				if i >= 1 {
					srcContext.lastToken = tokens[i-1]
				}

				if i < len(tokens)-1 {
					srcContext.nextToken = tokens[i+1]
				}

				if err := lint.CheckToken(srcContext); err != nil {
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
