// +build nrf,reelboard

package machine

const HasLowFrequencyCrystal = true

// LEDs on the reel board
const (
	LED        = LED_YELLOW
	LED_RED    = 11
	LED_GREEN  = 12
	LED_BLUE   = 41
	LED_YELLOW = 13
)

// User "a" button on the reel board
const (
	BUTTON = 7
)

// UART pins
const (
	UART_TX_PIN = 6
	UART_RX_PIN = 8
)

// I2C pins
const (
	SDA_PIN = 26
	SCL_PIN = 27
)
