package main

import (
	"github.com/Abbas-gheydi/radotp/pkgs/confs"
	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
	"github.com/Abbas-gheydi/radotp/pkgs/web"
)

func main() {

	confs.LoadConfigs()

	//database configuraion
	storage.Initialize()

	if confs.Cfg.Metrics.EnablePrometheusExporter {
		go monitoring.Start()
	}

	go web.StartRouter()

	rad.StartRadius()

}
