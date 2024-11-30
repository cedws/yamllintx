package lint

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"iter"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

var lintError = errors.New("lint error")

func newLintError(err error) error {
	return fmt.Errorf("%w: %w", lintError, err)
}

type tokenContext struct {
	lastToken    *token.Token
	currentToken *token.Token
	nextToken    *token.Token
}

type lineContext struct {
	currentLine       string
	currentLineNumber int
}

type Linter interface {
	CheckToken(tokenContext) iter.Seq[Problem]
	CheckLine(lineContext) iter.Seq[Problem]
}

type Chain []Linter

type Problem struct {
	Line   int
	Column int
	Error  error
}

// LintAll performs linting on the entire source code and returns an iterator of all errors found.
func LintAll(src []byte, linters ...Linter) iter.Seq[Problem] {
	tokens := lexer.Tokenize(string(src))

	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(src))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	seqFunc := func(yield func(Problem) bool) {
		for _, lint := range linters {
			for i := 0; i < len(lines); i++ {
				lineContext := lineContext{
					currentLine:       lines[i],
					currentLineNumber: i + 1,
				}

				for problem := range lint.CheckLine(lineContext) {
					if !yield(problem) {
						return
					}
				}
			}

			for i := 0; i < len(tokens); i++ {
				srcContext := tokenContext{
					currentToken: tokens[i],
				}

				if i >= 1 {
					srcContext.lastToken = tokens[i-1]
				}

				if i < len(tokens)-1 {
					srcContext.nextToken = tokens[i+1]
				}

				for problem := range lint.CheckToken(srcContext) {
					if !yield(problem) {
						return
					}
				}
			}
		}
	}

	return seqFunc
}

// Lint performs linting on the entire source code and returns the first error found.
func Lint(src []byte, linters ...Linter) *Problem {
	for err := range LintAll(src, linters...) {
		return &err
	}

	return nil
}
