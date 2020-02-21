package main

import (
	"flag"

	"github.com/UnityTech/nemesis/pkg/runner"
)

func main() {
	flag.Parse()
	audit := runner.NewAudit()
	audit.Setup()
	audit.Execute()
	audit.Report()
}
