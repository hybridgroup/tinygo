// +build stm32,stm32f0xx

package runtime

import (
	"device/arm"
)

const tickMicros = 1 // TODO

func init() {
	//machine.UART0.Configure(machine.UARTConfig{})
}

func putchar(c byte) {
	//machine.UART0.WriteByte(c)
}

func sleepTicks(d timeUnit) {
	// TODO: use a real timer here
	for i := 0; i < int(d/535); i++ {
		arm.Asm("")
	}
}

func ticks() timeUnit {
	return 0 // TODO
}
