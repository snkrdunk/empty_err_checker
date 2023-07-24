package main

import (
	"empty_err_checker"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(empty_err_checker.Analyzer) }
