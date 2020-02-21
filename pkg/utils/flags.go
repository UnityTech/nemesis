package utils

import "flag"

var (
	flagDebug = flag.Bool("debug", GetEnvBool("NEMESIS_DEBUG"), "Enable verbose output for debugging")
)
