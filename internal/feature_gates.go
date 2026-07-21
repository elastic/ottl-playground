package internal

import "go.opentelemetry.io/collector/featuregate"

func init() {
	enableFeatureGates()
}

func enableFeatureGates() {
	// OTTL lambda is used across executors.
	_ = featuregate.GlobalRegistry().Set("ottl.functions.enableLambda", true)
}
