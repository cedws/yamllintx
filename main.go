package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/cedws/yamllintx/lint"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

var ruleFactories = map[string]func() lint.Linter{
	"anchors": func() lint.Linter {
		return lint.Anchors(lint.AnchorOpts{})
	},
	"braces": func() lint.Linter {
		return lint.Braces{}
	},
	"brackets": func() lint.Linter {
		return lint.Brackets{}
	},
	"comments": func() lint.Linter {
		return lint.Comments{}
	},
	"hyphens": func() lint.Linter {
		return lint.Hyphens{}
	},
	"octal": func() lint.Linter {
		return lint.Octal{}
	},
	"trailing-spaces": func() lint.Linter {
		return lint.TrailingSpaces{}
	},
}

type config struct {
	YamlFiles []string            `yaml:"yaml-files"`
	Ignore    []string            `yaml:"ignore"`
	Rules     map[string]ast.Node `yaml:"rules"`
}

func unmarshalConfig(file string) config {
	config := config{
		YamlFiles: []string{
			"*.yaml",
			"*.yml",
			".yamllint",
		},
	}

	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {
	configFile := flag.String("config", "", "config file")
	flag.Parse()

	if configFile == nil || *configFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	config := unmarshalConfig(*configFile)
	var chain lint.Chain

	for rule := range config.Rules {
		if factory, ok := ruleFactories[rule]; ok {
			chain = append(chain, factory())
		}
	}

	files, err := lintableFiles(config, ".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Fprintln(os.Stderr, filepath.Clean(file))

		bytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		for err := range lint.LintAll(bytes, chain...) {
			fmt.Fprintf(os.Stderr, "  %d:%d\t%s\n", err.Line, err.Column, err.Error)
		}
	}
}

func lintableFiles(config config, dir string) ([]string, error) {
	var files []string

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		for _, pattern := range config.Ignore {
			match, err := doublestar.Match(pattern, path)
			if err != nil {
				return err
			}
			if match {
				return nil
			}
		}

		for _, pattern := range config.YamlFiles {
			match, err := doublestar.Match(pattern, d.Name())
			if err != nil {
				return err
			}
			if match {
				files = append(files, path)
			}
		}

		return nil
	}

	if err := filepath.WalkDir(dir, walkFunc); err != nil {
		return nil, err
	}

	return files, nil
}
