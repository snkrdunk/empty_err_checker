package main

import (
	"github.com/snkrdunk/empty_err_checker"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(empty_err_checker.Analyzer) }
