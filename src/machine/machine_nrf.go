// +build nrf

package machine

import (
	"device/arm"
	"device/nrf"
	"errors"
)

type GPIOMode uint8

const (
	GPIO_INPUT          = (nrf.GPIO_PIN_CNF_DIR_Input << nrf.GPIO_PIN_CNF_DIR_Pos) | (nrf.GPIO_PIN_CNF_INPUT_Connect << nrf.GPIO_PIN_CNF_INPUT_Pos)
	GPIO_INPUT_PULLUP   = GPIO_INPUT | (nrf.GPIO_PIN_CNF_PULL_Pullup << nrf.GPIO_PIN_CNF_PULL_Pos)
	GPIO_INPUT_PULLDOWN = GPIO_INPUT | (nrf.GPIO_PIN_CNF_PULL_Pulldown << nrf.GPIO_PIN_CNF_PULL_Pos)
	GPIO_OUTPUT         = (nrf.GPIO_PIN_CNF_DIR_Output << nrf.GPIO_PIN_CNF_DIR_Pos) | (nrf.GPIO_PIN_CNF_INPUT_Disconnect << nrf.GPIO_PIN_CNF_INPUT_Pos)
)

// Configure this pin with the given configuration.
func (p GPIO) Configure(config GPIOConfig) {
	cfg := config.Mode | nrf.GPIO_PIN_CNF_DRIVE_S0S1 | nrf.GPIO_PIN_CNF_SENSE_Disabled
	nrf.P0.PIN_CNF[p.Pin] = nrf.RegValue(cfg)
}

// Set the pin to high or low.
// Warning: only use this on an output pin!
func (p GPIO) Set(high bool) {
	if high {
		nrf.P0.OUTSET = 1 << p.Pin
	} else {
		nrf.P0.OUTCLR = 1 << p.Pin
	}
}

// Get returns the current value of a GPIO pin.
func (p GPIO) Get() bool {
	return (nrf.P0.IN>>p.Pin)&1 != 0
}

// UART
type UARTConfig struct {
	Baudrate uint32
}

type UART struct {
}

var (
	// UART0 is the hardware serial port on the NRF.
	UART0 = &UART{}
)

// Configure the UART.
func (uart UART) Configure(config UARTConfig) {
	// Default baud rate to 115200.
	if config.Baudrate == 0 {
		config.Baudrate = 115200
	}

	uart.SetBaudRate(config.Baudrate)

	// TODO: Set TX and RX pins from config.
	nrf.UART0.PSELTXD = 6 // pin 6 for NRF52840-DK
	nrf.UART0.PSELRXD = 8 // pin 8 for NRF52840-DK

	nrf.UART0.ENABLE = nrf.UART_ENABLE_ENABLE_Enabled
	nrf.UART0.TASKS_STARTTX = 1
	nrf.UART0.TASKS_STARTRX = 1
	nrf.UART0.INTENSET = nrf.UART_INTENSET_RXDRDY_Msk

	// Enable RX IRQ.
	arm.EnableIRQ(nrf.IRQ_UARTE0_UART0)
}

// SetBaudRate sets the communication speed for the UART.
func (uart UART) SetBaudRate(br uint32) {
	switch br {
	case 1200:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud1200
	case 2400:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud2400
	case 4800:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud4800
	case 9600:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud9600
	case 14400:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud14400
	case 19200:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud19200
	case 28800:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud28800
	case 38400:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud38400
	case 57600:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud57600
	case 76800:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud76800
	case 115200:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud115200
	case 230400:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud230400
	case 250000:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud250000
	case 460800:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud460800
	case 921600:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud921600
	case 1000000:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud1M
	default:
		nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud115200
	}
}

// Close the UART.
func (uart UART) Close() error {
	return nil
}

// Read from the RX buffer.
func (uart UART) Read(data []byte) (n int, err error) {
	if len(data) > RXBufferSize {
		return 0, errors.New("Read buffer cannot be larger than RXBuffer")
	}

	// check if RX buffer is empty
	size := BufferUsed()
	if size == 0 {
		return 0, nil
	}

	// only read number of bytes used from buffer
	for i := 0; uint8(i) < size; i++ {
		data[i] = byte(bufferGet())
	}

	return len(data), nil
}

// Write data to the UART.
func (uart UART) Write(data []byte) (n int, err error) {
	for _, v := range data {
		uart.WriteByte(v)
	}
	return len(data), nil
}

// ReadByte reads a single byte from the RX buffer.
func (uart UART) ReadByte() (byte, error) {
	return byte(bufferGet()), nil
}

// WriteByte writes a byte of data to the UART.
func (uart UART) WriteByte(c byte) error {
	nrf.UART0.EVENTS_TXDRDY = 0
	nrf.UART0.TXD = nrf.RegValue(c)
	for nrf.UART0.EVENTS_TXDRDY == 0 {
	}
	return nil
}

type __volatile byte

const RXBufferSize = 64

// Minimal ring buffer implementation inspired by post at
// https://www.embeddedrelated.com/showthread/comp.arch.embedded/77084-1.php
var rxbuffer [RXBufferSize]__volatile
var head __volatile
var tail __volatile

func BufferUsed() uint8 {
	return uint8(bufferUsed())
}

func bufferUsed() __volatile { return head - tail }
func bufferPut(val __volatile) {
	if bufferUsed() != RXBufferSize {
		head++
		rxbuffer[head%RXBufferSize] = val
	}
}
func bufferGet() __volatile {
	if bufferUsed() != 0 {
		tail++
		return rxbuffer[tail%RXBufferSize]
	}
	return 0
}

//go:export UARTE0_UART0_IRQHandler
func handleUART0() {
	if nrf.UART0.EVENTS_RXDRDY == 1 {
		bufferPut(__volatile(nrf.UART0.RXD))
		nrf.UART0.EVENTS_RXDRDY = 0x0
	}
}
