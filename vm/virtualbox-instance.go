package vm

import (
	"log"

	vbox "github.com/rmxymh/go-virtualbox"
)

type VirtualBoxInstance struct {
	Name string

	defaultBootOrder []string
	nextBootOrder    []string
	changeBootOrder  bool
}

func NewVirtualBoxInstance(name string) *VirtualBoxInstance {
	return &VirtualBoxInstance{
		Name:             name,
		defaultBootOrder: []string{"disk", "net"},
	}
}

func (instance *VirtualBoxInstance) IsRunning() bool {
	machine, err := vbox.GetMachine(instance.Name)

	if err == nil && machine.State == vbox.Running {
		return true
	}
	return false
}

func (instance *VirtualBoxInstance) SetBootDevice(dev string) {
	machine, err := vbox.GetMachine(instance.Name)

	if err != nil {
		log.Fatalf("    Instance: Failed to set BootDevice to VM %s: %s", instance.Name, err.Error())
		return
	}

	newBootOrder := []string{dev}
	for _, d := range machine.BootOrder {
		if d != dev {
			newBootOrder = append(newBootOrder, d)
		}
	}

	instance.nextBootOrder = newBootOrder
	instance.changeBootOrder = true
}

func (instance *VirtualBoxInstance) PowerOff() {
	machine, err := vbox.GetMachine(instance.Name)

	if err != nil {
		log.Fatalf("    Instance: Failed to find VM %s and power off it: %s", instance.Name, err.Error())
		return
	}

	machine.Poweroff()
}

func (instance *VirtualBoxInstance) ACPIOff() {
	machine, err := vbox.GetMachine(instance.Name)

	if err != nil {
		log.Fatalf("    Instance: Failed to find VM %s and power off it: %s", instance.Name, err.Error())
		return
	}

	machine.Stop()
}

func (instance *VirtualBoxInstance) PowerOn() {
	machine, err := vbox.GetMachine(instance.Name)

	if err != nil {
		log.Fatalf("    Instance: Failed to find VM %s and power on it: %s", instance.Name, err.Error())
		return
	}

	if instance.changeBootOrder {
		machine.BootOrder = instance.nextBootOrder
		machine.Modify()
		instance.nextBootOrder = make([]string, 4)
		instance.changeBootOrder = false
		log.Println("Current Boot Order = ", machine.BootOrder)
	} else {
		machine.BootOrder = instance.defaultBootOrder
		machine.Modify()
		instance.changeBootOrder = false
		log.Println("Current Boot Order = ", machine.BootOrder)
	}

	machine.Start()
}

func (instance *VirtualBoxInstance) Reset() {
	machine, err := vbox.GetMachine(instance.Name)

	if err != nil {
		log.Fatalf("    Instance: Failed to find VM %s and power on it: %s", instance.Name, err.Error())
		return
	}

	machine.Reset()
}

func (instance *VirtualBoxInstance) NICInitialize() {
	machine, err := vbox.GetMachine(instance.Name)

	if err != nil {
		log.Fatalf("    Instance: Failed to find VM %s and power on it: %s", instance.Name, err.Error())
		return
	}

	nic := machine.NICs[0]
	// force NIC1 to internal network so that we can deploy it via PXE.
	nic.Network = vbox.NICNetInternal
	machine.SetNIC(1, nic)
}
