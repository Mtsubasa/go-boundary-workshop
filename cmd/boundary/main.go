package main

import (
	"github.com/Mtsubasa/go-boundary-workshop/boundary"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(boundary.Analyzer)
}
