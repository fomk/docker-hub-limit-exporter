package config

import (
	"flag"
	"os"
)

var (
	// Default values
	defaultListenAddress = getEnv("LISTEN_ADDRESS", ":9032")
	defaultMetricsPath   = getEnv("TELEMETRY_PATH", "/metrics")
	defaultUser          = getEnv("DOCKER_HUB_USER", "")
	defaultPass          = getEnv("DOCKER_HUB_PASS", "")

	// CLI flags
	ListenAddr = flag.String("web.listen-address",
		defaultListenAddress,
		"An address to listen on for web interface and telemetry. Can be overwritten by LISTEN_ADDRESS environment variable.")

	MetricsPath = flag.String("web.telemetry-path",
		defaultMetricsPath,
		"A path under which to expose metrics. Can be overwritten by TELEMETRY_PATH environment variable.")

	HubUser = flag.String("hub.user",
		defaultUser,
		"Docker hub user, anonymous used when empty. Can be overwritten by DOCKER_HUB_USER environment variable.")

	HubPass = flag.String("hub.pass",
		defaultPass,
		"Docker hub password. Can be overwritten by DOCKER_HUB_PASS environment variable.")
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}