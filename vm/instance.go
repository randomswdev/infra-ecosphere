package vm

import (
	"log"
)

type InstanceType string

const (
	Fake       InstanceType = "Fake"
	VirtualBox InstanceType = "VirtualBox"
	HyperV     InstanceType = "Hyper-V"
)

const (
	BOOT_DEVICE_PXE    = "net"
	BOOT_DEVICE_DISK   = "disk"
	BOOT_DEVICE_CD_DVD = "dvd"
	BOOT_DEVICE_FLOPPY = "floppy"
)

type Instance interface {
	IsRunning() bool
	SetBootDevice(dev string)
	PowerOff()
	ACPIOff()
	PowerOn()
	Reset()
	NICInitialize()
}

var instances map[string]Instance

func init() {
	instances = make(map[string]Instance)
}

func AddInstnace(name string, instanceType InstanceType) Instance {
	var newInstance Instance

	switch instanceType {
	case VirtualBox:
		newInstance = NewVirtualBoxInstance(name)
		break
	case HyperV:
		newInstance = NewHypervInstance(name)
		break
	case Fake:
		newInstance = NewFakeInstance(name)
		break
	default:
		log.Println("Unknown instance type  ", instanceType)
		newInstance = NewFakeInstance(name)
		break
	}

	instances[name] = newInstance
	newInstance.NICInitialize()
	log.Println("Add instance ", name)

	return newInstance
}

func DeleteInstance(name string) {
	_, ok := instances[name]
	if ok {
		delete(instances, name)
	}
	log.Println("Remove instance ", name)
}

func GetInstance(name string) (instance Instance, ok bool) {
	instance, ok = instances[name]
	return instance, ok
}
