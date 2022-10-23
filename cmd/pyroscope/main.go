package main

import (
	"runtime"

	"github.com/Abbas-gheydi/radotp/pkgs/confs"
	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
	"github.com/Abbas-gheydi/radotp/pkgs/web"
	"github.com/pyroscope-io/client/pyroscope"
)

func main() {
	// These 2 lines are only required if you're using mutex or block profiling
	// Read the explanation below for how to set these rates:
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	pyroscope.Start(pyroscope.Config{
		ApplicationName: "radotp",

		// replace this with the address of pyroscope server
		ServerAddress: "http://127.0.0.1:4040",

		// you can disable logging by setting this to nil
		Logger: pyroscope.StandardLogger,

		// optionally, if authentication is enabled, specify the API key:
		// AuthToken: os.Getenv("PYROSCOPE_AUTH_TOKEN"),

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})

	// your code goes here
	confs.LoadConfigs()

	//database configuraion
	storage.Initialize()

	if confs.Cfg.Metrics.EnablePrometheusExporter {
		go monitoring.Start()
	}

	go web.StartRouter()

	rad.StartRadius()
}
