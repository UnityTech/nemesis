package report

import (
	"flag"

	"github.com/UnityTech/nemesis/pkg/utils"
)

var (
	flagReportOnlyFailures = flag.Bool("reports.only-failures", utils.GetEnvBool("NEMESIS_ONLY_FAILURES"), "Limit output of controls to only failed controls")
)
