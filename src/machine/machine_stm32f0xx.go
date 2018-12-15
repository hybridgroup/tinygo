// +build stm32,stm32f0xx

package machine

import "device/stm32"

// Peripheral abstraction layer for the stm32f0xx.

const (
	GPIO_INPUT        = 0 // Input mode
	GPIO_OUTPUT       = 1 // Output mode
	GPIO_ALT          = 2 // Alt mode
	GPIO_INPUT_ANALOG = 3 // Input analog mode

	GPIO_OUTPUT_TYPE_Pos        = 2
	GPIO_OUTPUT_TYPE_PUSH_PULL  = 0 // Output type push/pull
	GPIO_OUTPUT_TYPE_OPEN_DRAIN = 1 // Output type open drain

	GPIO_OUTPUT_Speed_Pos = 4
	GPIO_OUTPUT_2MHz      = 0 // Output speed 2MHz
	GPIO_OUTPUT_10MHz     = 1 // Output speed 10MHz
	GPIO_OUTPUT_50MHz     = 3 // Output speed 50MHz

	GPIO_PullupDown_Pos = 6
	GPIO_FLOATING       = 0 // Floating mode
	GPIO_PULL_UP        = 1 // Pull up
	GPIO_PULL_DOWN      = 2 // Pull down
	GPIO_RESERVED       = 3 // reserved
)

func (p GPIO) getPort() *stm32.GPIO_Type {
	switch p.Pin / 16 {
	case 0:
		return stm32.GPIOA
	case 1:
		return stm32.GPIOB
	case 2:
		return stm32.GPIOC
	case 3:
		return stm32.GPIOD
	case 4:
		panic("machine: unknown port")
	case 5:
		return stm32.GPIOF
	default:
		panic("machine: unknown port")
	}
}

func (p GPIO) enableClock() {
	switch p.Pin / 16 {
	case 0:
		stm32.RCC.AHBENR |= stm32.RCC_AHBENR_IOPAEN
	case 1:
		stm32.RCC.AHBENR |= stm32.RCC_AHBENR_IOPBEN
	case 2:
		stm32.RCC.AHBENR |= stm32.RCC_AHBENR_IOPCEN
	case 3:
		stm32.RCC.AHBENR |= stm32.RCC_AHBENR_IOPDEN
	case 4:
		panic("machine: unknown port")
	case 5:
		stm32.RCC.AHBENR |= stm32.RCC_AHBENR_IOPFEN
	default:
		panic("machine: unknown port")
	}
}

// Configure this pin with the given configuration.
func (p GPIO) Configure(config GPIOConfig) {
	p.enableClock()
	port := p.getPort()
	pos := p.Pin % 8 * 2

	mode := config.Mode & 0xf
	if mode == GPIO_OUTPUT {
		// Configure the GPIO output pin type
		port.OTYPER = stm32.RegValue((uint32(port.OTYPER) &^ (0xf << p.Pin)) |
			(uint32((config.Mode>>GPIO_OUTPUT_TYPE_Pos)&0xf) << p.Pin))

		// Configure the GPIO output pin speed
		port.OSPEEDR = stm32.RegValue((uint32(port.OSPEEDR) &^ (0xf << p.Pin)) |
			(uint32((config.Mode>>GPIO_OUTPUT_Speed_Pos)&0xf) << p.Pin))

		// Configure the GPIO output pin pull-up/pull-down resistor
		port.PUPDR = stm32.RegValue((uint32(port.PUPDR) &^ (0xf << p.Pin)) |
			(uint32((config.Mode>>GPIO_PullupDown_Pos)&0xf) << p.Pin))
	}

	// Configure the GPIO pin mode
	port.MODER = stm32.RegValue((uint32(port.MODER) &^ (0xf << pos)) | (uint32(mode) << pos))
}

// Set the pin to high or low.
// Warning: only use this on an output pin!
func (p GPIO) Set(high bool) {
	port := p.getPort()
	pin := p.Pin % 16
	if high {
		port.BSRR = 1 << pin
	} else {
		port.BSRR = 1 << (pin + 16)
	}
}
