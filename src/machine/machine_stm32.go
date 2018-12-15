// +build stm32

package machine

// Peripheral abstraction layer for the stm32.

type GPIOMode uint8

const SystemClockSpeed = 72000000

const (
	portA = iota * 16
	portB
	portC
	portD
	portE
	portF
	portG
)
