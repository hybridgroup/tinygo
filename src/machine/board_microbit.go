// +build nrf,microbit

package machine

import "device/nrf"

// LEDs on the BBC:Microbit (nRF51 dev board)
const (
	LED  = LED1
	LED1 = 1
)

// Buttons on the BBC:Microbit (nRF51 dev board)
const (
	BUTTON  = BUTTON1
	BUTTON1 = 5
	BUTTON2 = 11
)

func SetLEDS(on bool) {
	nrf.GPIO.DIRSET = nrf.RegValue(0xfff0)
	// v = 0                       // All pins to GND so no LED lid
	// | MASK_LED_ROW_0
	// | MASK_LED_ROW_1
	// | MASK_LED_ROW_2
	// | MASK_LED_COL_0
	// | MASK_LED_COL_1
	// | MASK_LED_COL_2
	// | MASK_LED_COL_3
	// | MASK_LED_COL_4
	// | MASK_LED_COL_5
	// | MASK_LED_COL_6
	// | MASK_LED_COL_7
	// | MASK_LED_COL_8
	if on {
		nrf.GPIO.OUTSET = nrf.RegValue(0xfff0)
	} else {
		nrf.GPIO.OUTCLR = nrf.RegValue(0xfff0)
	}
}
