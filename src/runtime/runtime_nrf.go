// +build nrf

package runtime

import (
	"device/arm"
	"device/nrf"
)

type timeUnit int64

const tickMicros = 1024 * 32

//go:linkname systemInit SystemInit
func systemInit()

//go:export Reset_Handler
func handleReset() {
	systemInit()
	main()
}

func init() {
	initUART()
	initLFCLK()
	initRTC()
}

func initUART() {
	nrf.UART0.ENABLE = nrf.UART_ENABLE_ENABLE_Enabled
	nrf.UART0.BAUDRATE = nrf.UART_BAUDRATE_BAUDRATE_Baud9600
	nrf.UART0.TASKS_STARTTX = 1
	nrf.UART0.PSELTXD = 24 // pin 6 for NRF52840-DK
	nrf.UART0.PSELRXD = 25 // pin 6 for NRF52840-DK
}

func initLFCLK() {
	nrf.CLOCK.LFCLKSRC = nrf.CLOCK_LFCLKSTAT_SRC_Xtal
	nrf.CLOCK.TASKS_LFCLKSTART = 1
	for nrf.CLOCK.EVENTS_LFCLKSTARTED == 0 {
	}
	nrf.CLOCK.EVENTS_LFCLKSTARTED = 0
}

func initRTC() {
	nrf.RTC0.TASKS_START = 1
	// TODO: set priority
	arm.EnableIRQ(nrf.IRQ_RTC0)
}

func putchar(c byte) {
	nrf.UART0.EVENTS_TXDRDY = 0
	nrf.UART0.TXD = nrf.RegValue(c)
	for nrf.UART0.EVENTS_TXDRDY == 0 {
	}
}

func sleepTicks(d timeUnit) {
	for d != 0 {
		ticks()                       // update timestamp
		ticks := uint32(d) & 0x7fffff // 23 bits (to be on the safe side)
		rtc_sleep(ticks)              // TODO: not accurate (must be d / 30.5175...)
		d -= timeUnit(ticks)
	}
}

var (
	timestamp      timeUnit // nanoseconds since boottime
	rtcLastCounter uint32   // 24 bits ticks
)

// Monotonically increasing numer of ticks since start.
//
// Note: very long pauses between measurements (more than 8 minutes) may
// overflow the counter, leading to incorrect results. This might be fixed by
// handling the overflow event.
func ticks() timeUnit {
	rtcCounter := uint32(nrf.RTC0.COUNTER)
	offset := (rtcCounter - rtcLastCounter) % 0xffffff // change since last measurement
	rtcLastCounter = rtcCounter
	timestamp += timeUnit(offset) // TODO: not precise
	return timestamp
}

//go:volatile
type isrFlag bool

var rtc_wakeup isrFlag

func rtc_sleep(ticks uint32) {
	nrf.RTC0.INTENSET = nrf.RTC_INTENSET_COMPARE0
	rtc_wakeup = false
	if ticks == 1 {
		// Race condition (even in hardware) at ticks == 1.
		// TODO: fix this in a better way by detecting it, like the manual
		// describes.
		ticks = 2
	}
	nrf.RTC0.CC[0] = (nrf.RTC0.COUNTER + nrf.RegValue(ticks)) & 0x00ffffff
	for !rtc_wakeup {
		arm.Asm("wfi")
	}
}

//go:export RTC0_IRQHandler
func handleRTC0() {
	nrf.RTC0.INTENCLR = nrf.RTC_INTENSET_COMPARE0
	nrf.RTC0.EVENTS_COMPARE[0] = 0
	rtc_wakeup = true
}
