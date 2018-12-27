// Hand created file. DO NOT DELETE.
// STM32F103XX bitfield definitions that are not auto-generated by gen-device-svd.py

// +build stm32,stm32f103xx

package stm32

const (
	// Flash Access Control Register flag values.
	FLASH_ACR_LATENCY_0 = 0x00000001
	FLASH_ACR_LATENCY_1 = 0x00000002
	FLASH_ACR_LATENCY_2 = 0x00000004

	// Reset and Clock Control Control Register flag values.

	// System Clock source
	RCC_CFGR_SW_HSI = 0
	RCC_CFGR_SW_HSE = 1
	RCC_CFGR_SW_PLL = 2

	// Flags for when System Clock source is set.
	RCC_CFGR_SWS_HSI = 0x00000000
	RCC_CFGR_SWS_HSE = 0x00000004
	RCC_CFGR_SWS_PLL = 0x00000008

	// Sets PCLK1
	RCC_CFGR_PPRE1_DIV_NONE = 0x00000000
	RCC_CFGR_PPRE1_DIV_2    = 0x00000400
	RCC_CFGR_PPRE1_DIV_4    = 0x00000500
	RCC_CFGR_PPRE1_DIV_8    = 0x00000600
	RCC_CFGR_PPRE1_DIV_16   = 0x00000700

	// Sets PCLK2
	RCC_CFGR_PPRE2_DIV_NONE = 0x00000000
	RCC_CFGR_PPRE2_DIV_2    = 0x00002000
	RCC_CFGR_PPRE2_DIV_4    = 0x00002800
	RCC_CFGR_PPRE2_DIV_8    = 0x00003000
	RCC_CFGR_PPRE2_DIV_16   = 0x00003800

	// Sets PLL multiplier
	RCC_CFGR_PLLMUL_2  = 0x00000000
	RCC_CFGR_PLLMUL_3  = 0x00040000
	RCC_CFGR_PLLMUL_4  = 0x00080000
	RCC_CFGR_PLLMUL_5  = 0x000C0000
	RCC_CFGR_PLLMUL_6  = 0x00100000
	RCC_CFGR_PLLMUL_7  = 0x00140000
	RCC_CFGR_PLLMUL_8  = 0x00180000
	RCC_CFGR_PLLMUL_9  = 0x001C0000
	RCC_CFGR_PLLMUL_10 = 0x00200000
	RCC_CFGR_PLLMUL_11 = 0x00240000
	RCC_CFGR_PLLMUL_12 = 0x00280000
	RCC_CFGR_PLLMUL_13 = 0x002C0000
	RCC_CFGR_PLLMUL_14 = 0x00300000
	RCC_CFGR_PLLMUL_15 = 0x00340000
	RCC_CFGR_PLLMUL_16 = 0x00380000

	// RTC clock source
	RCC_RTCCLKSource_LSE        = 0x00000100
	RCC_RTCCLKSource_LSI        = 0x00000200
	RCC_RTCCLKSource_HSE_Div128 = 0x00000300
)
