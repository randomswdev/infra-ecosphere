package vm

import (
	"bytes"
	"log"
	"text/template"

	"github.com/hashicorp/packer/common/powershell"
	"github.com/hashicorp/packer/common/powershell/hyperv"
)

type HypervInstance struct {
	Name string

	defaultBootOrder []string
	nextBootOrder    []string
	changeBootOrder  bool
}

const powerShellTemplateGen1 = `
{{$out := .}}
param([string]$vmName)
Hyper-V\Set-VMBios -VMName $vmName -StartupOrder @({{ range $index, $device := .Devices}}{{if $index}},{{end}}{{index $out.DeviceMap 0 $device}}{{end}})
`

const powerShellTemplateGen2 = `
{{$out := .}}
param([string]$vmName)
{{ range $index, $device := .Devices}}
{{index $out.GettersMap $device}}
{{end}}
Hyper-V\Set-VMFirmware -VMName $vmName -BootOrder @({{ range $index, $device := .Devices}}{{if $index}},{{end}}{{index $out.DeviceMap 1 $device}}{{end}})
`

type templateParameters struct {
	Devices    []string
	DeviceMap  []map[string]string
	GettersMap map[string]string
}

var gettersMap map[string]string
var deviceMap []map[string]string

func init() {
	deviceMap = make([]map[string]string, 2)

	deviceMap[0] = make(map[string]string)
	deviceMap[0][BOOT_DEVICE_PXE] = "\"LegacyNetworkAdapter\""
	deviceMap[0][BOOT_DEVICE_DISK] = "\"IDE\""
	deviceMap[0][BOOT_DEVICE_CD_DVD] = "\"CD\""
	deviceMap[0][BOOT_DEVICE_FLOPPY] = "\"Floppy\""

	deviceMap[1] = make(map[string]string)
	deviceMap[1][BOOT_DEVICE_PXE] = "$MyNIC"
	deviceMap[1][BOOT_DEVICE_DISK] = "$MyHD"
	deviceMap[1][BOOT_DEVICE_CD_DVD] = "$MyDVD"
	deviceMap[1][BOOT_DEVICE_FLOPPY] = ""

	gettersMap = make(map[string]string)
	gettersMap[BOOT_DEVICE_PXE] = "$MyNIC = Get-VMNetworkAdapter $vmName"
	gettersMap[BOOT_DEVICE_DISK] = "$MyHD = Get-VMHardDiskDrive $vmName"
	gettersMap[BOOT_DEVICE_CD_DVD] = "$MyDVD = Get-VMDvdDrive $vmName"
	gettersMap[BOOT_DEVICE_FLOPPY] = ""
}

func NewHypervInstance(name string) *HypervInstance {
	return &HypervInstance{
		Name:             name,
		defaultBootOrder: []string{"disk", "net"},
	}
}

func (instance *HypervInstance) IsRunning() bool {
	running, _ := hyperv.IsRunning(instance.Name)
	return running
}

func (instance *HypervInstance) SetBootDevice(dev string) {
	instance.nextBootOrder = []string{dev}
	instance.changeBootOrder = true
}

func (instance *HypervInstance) PowerOff() {
	hyperv.ShutDown(instance.Name)
}

func (instance *HypervInstance) ACPIOff() {
	hyperv.TurnOff(instance.Name)
}

func (instance *HypervInstance) PowerOn() {
	generation, err := hyperv.GetVirtualMachineGeneration(instance.Name)
	if err != nil {
		log.Fatalf("    Instance: Failed to get retrieve generation for VM %s: %s", instance.Name, err.Error())
		return
	}

	bootOrder := instance.defaultBootOrder

	if instance.changeBootOrder {
		bootOrder = instance.nextBootOrder

		instance.nextBootOrder = make([]string, 4)
		instance.changeBootOrder = false
	}

	log.Println("Current Boot Order = ", bootOrder)

	parameters := templateParameters{
		Devices:    bootOrder,
		DeviceMap:  deviceMap,
		GettersMap: gettersMap,
	}

	powerShellTemplate := powerShellTemplateGen2
	if generation < 2 {
		powerShellTemplate = powerShellTemplateGen1
	}

	var buffer bytes.Buffer
	compiledPowerShellTemplate, err := template.New("script").Parse(powerShellTemplate)
	if err != nil {
		log.Fatalf("    Instance: Failed to parse boot command for VM %s: %s", instance.Name, err.Error())
		return
	}
	if err := compiledPowerShellTemplate.Execute(&buffer, parameters); err != nil {
		log.Fatalf("    Instance: Failed to generate boot command for VM %s: %s", instance.Name, err.Error())
		return
	}

	var ps powershell.PowerShellCmd
	err = ps.Run(buffer.String(), instance.Name)

	hyperv.StartVirtualMachine(instance.Name)
}

func (instance *HypervInstance) Reset() {
	hyperv.RestartVirtualMachine(instance.Name)
}

func (instance *HypervInstance) NICInitialize() {
	// Nothing to do here
	return
}
