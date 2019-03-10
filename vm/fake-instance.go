package vm

type FakeInstance struct {
}

func NewFakeInstance(_ string) *FakeInstance {
	return &FakeInstance{}
}

func (instance *FakeInstance) IsRunning() bool {
	return true
}

func (instance *FakeInstance) SetBootDevice(dev string) {
	return
}

func (instance *FakeInstance) PowerOff() {
	return
}

func (instance *FakeInstance) ACPIOff() {
	return
}

func (instance *FakeInstance) PowerOn() {
	return
}

func (instance *FakeInstance) Reset() {
	return
}

func (instance *FakeInstance) NICInitialize() {
	return
}
