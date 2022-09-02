package main

import (
	"fmt"

	"github.com/Abbas-gheydi/radotp/pkgs/confs"
	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
	"github.com/Abbas-gheydi/radotp/pkgs/web"
)

func art() {

	fmt.Print(`	 
                                                                  
    _/_/_/      _/_/    _/_/_/      _/_/    _/_/_/_/_/  _/_/_/    
   _/    _/  _/    _/  _/    _/  _/    _/      _/      _/    _/   
  _/_/_/    _/_/_/_/  _/    _/  _/    _/      _/      _/_/_/      
 _/    _/  _/    _/  _/    _/  _/    _/      _/      _/           
_/    _/  _/    _/  _/_/_/      _/_/        _/      _/            
                                                                  
                                                                  
                                                              
`)
}

func main() {
	confs.LoadConfigs()

	//art()

	//database configuraion
	storage.Initialize()

	if confs.Cfg.Metrics.EnablePrometheusExporter {
		go monitoring.Start()
	}

	go web.StartRouter()

	rad.StartRadius()

}
