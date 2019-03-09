package main

import (
	"infra-ecosphere/ipmi"
	"infra-ecosphere/utils"
	"infra-ecosphere/web"
)

func main() {
	utils.LoadConfig("infra-ecosphere.cfg")
	go ipmi.IPMIServerServiceRun()
	web.WebAPIServiceRun()
}
