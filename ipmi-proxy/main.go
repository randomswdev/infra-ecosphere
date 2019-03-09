package main

import (
	"infra-ecosphere/ipmi"
	"infra-ecosphere/utils"
)

func main() {
	config := utils.LoadConfig("infra-ecosphere.cfg")
	EcospherePort = config.WebAPIPort
	ipmi.IPMI_CHASSIS_SET_BOOT_OPTION_SetHandler(ipmi.BOOT_FLAG, SetBootDevice)
	ipmi.IPMI_CHASSIS_SetHandler(ipmi.IPMI_CMD_GET_CHASSIS_STATUS, HandleIPMIGetChassisStatus)
	ipmi.IPMI_CHASSIS_SetHandler(ipmi.IPMI_CMD_CHASSIS_CONTROL, HandleIPMIChassisControl)
	ipmi.IPMIServerServiceRun()
}
