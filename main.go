package main

import (
	"flag"

	"github.com/Unity-Technologies/nemesis/pkg/runner"
)

func main() {
	flag.Parse()
	audit := runner.NewAudit()
	audit.Setup()
	audit.Execute()
	audit.Report()
}
