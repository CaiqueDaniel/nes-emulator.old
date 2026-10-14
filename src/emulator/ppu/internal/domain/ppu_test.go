package domain

import (
	"testing"
)

func TestNewPPU(t *testing.T) {
	sut := NewPPU()

	if sut == nil {
		t.Fatal("expected NewPPU to return a non-nil instance")
	}

	if sut.GetCurrentScanline() != 0 {
		t.Errorf("expected initial scanline to be 0, got %d", sut.GetCurrentScanline())
	}

	if sut.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected initial pixel to be 0, got %d", sut.GetCurrentScanlinePixel())
	}

	if sut.GetCurrentDots() != 0 {
		t.Errorf("expected initial dots to be 0, got %d", sut.GetCurrentDots())
	}

	if sut.GetV() != 0 {
		t.Errorf("expected initial v to be 0, got %d", sut.GetV())
	}

	if sut.GetT() != 0 {
		t.Errorf("expected initial t to be 0, got %d", sut.GetT())
	}

	if sut.GetX() != 0 {
		t.Errorf("expected initial x to be 0, got %d", sut.GetX())
	}

	if sut.IsWSet() {
		t.Errorf("expected initial w to be false, got true")
	}

	if sut.IsRenderingEnabled() {
		t.Errorf("expected enableRender to be false, got true")
	}
}

func TestPPU_AdvanceToNextScanlinePixel(t *testing.T) {
	sut := NewPPU()

	sut.AdvanceToNextScanlinePixel()
	if sut.GetCurrentScanlinePixel() != 1 {
		t.Errorf("expected pixel to be 1, got %d", sut.GetCurrentScanlinePixel())
	}

	for i := 0; i < 254; i++ {
		sut.AdvanceToNextScanlinePixel()
	}
	if sut.GetCurrentScanlinePixel() != 255 {
		t.Errorf("expected pixel to be 255, got %d", sut.GetCurrentScanlinePixel())
	}

	// Test overflow from 255 to 0 (uint8)
	sut.AdvanceToNextScanlinePixel()
	if sut.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected pixel to overflow to 0, got %d", sut.GetCurrentScanlinePixel())
	}
}

func TestPPU_IncreaseDots(t *testing.T) {
	sut := NewPPU()

	for expectedDot := uint16(1); expectedDot < MAX_DOTS_PER_LINE; expectedDot++ {
		sut.IncreaseDots()
		if sut.GetCurrentDots() != expectedDot {
			t.Fatalf("expected dots to be %d, got %d", expectedDot, sut.GetCurrentDots())
		}
	}

	// At MAX_DOTS_PER_LINE - 1 (335), increasing dots wraps around to 0
	sut.IncreaseDots()
	if sut.GetCurrentDots() != 0 {
		t.Errorf("expected dots to wrap around to 0, got %d", sut.GetCurrentDots())
	}
}

func TestPPU_AdvanceToNextScanline(t *testing.T) {
	t.Run("should not advance scanline if dots != MAX_DOTS_PER_LINE-1", func(t *testing.T) {
		sut := NewPPU()

		sut.AdvanceToNextScanline()
		if sut.GetCurrentScanline() != 0 {
			t.Errorf("expected scanline to remain 0 when dots=0, got %d", sut.GetCurrentScanline())
		}

		sut.dots = 334
		sut.AdvanceToNextScanline()
		if sut.GetCurrentScanline() != 0 {
			t.Errorf("expected scanline to remain 0 when dots=334, got %d", sut.GetCurrentScanline())
		}
	})

	t.Run("should advance scanline when dots == MAX_DOTS_PER_LINE-1", func(t *testing.T) {
		sut := NewPPU()
		sut.dots = MAX_DOTS_PER_LINE - 1

		sut.AdvanceToNextScanline()
		if sut.GetCurrentScanline() != 1 {
			t.Errorf("expected scanline to advance to 1, got %d", sut.GetCurrentScanline())
		}
	})

	t.Run("should wrap scanline to 0 after MAX_FRAME_SCANLINE", func(t *testing.T) {
		sut := NewPPU()
		sut.dots = MAX_DOTS_PER_LINE - 1
		sut.scanline = MAX_FRAME_SCANLINE

		sut.AdvanceToNextScanline()
		if sut.GetCurrentScanline() != 0 {
			t.Errorf("expected scanline to wrap to 0, got %d", sut.GetCurrentScanline())
		}
	})
}

func TestPPU_SetV_GetV(t *testing.T) {
	sut := NewPPU()

	testCases := []struct {
		input    []byte
		expected uint16
	}{
		{input: []byte{0x00, 0x00}, expected: 0x0000},
		{input: []byte{0x12, 0x34}, expected: 0x1234},
		{input: []byte{0x7F, 0xFF}, expected: 16383},
		{input: []byte{0x80, 0x00}, expected: 0},
		{input: []byte{0xFF, 0xFF}, expected: 16383},
	}

	for _, tc := range testCases {
		for _, b := range tc.input {
			sut.AppendOnV(b)
		}

		if sut.GetV() != tc.expected {
			t.Errorf("SetV(%d): expected %d, got %d", tc.input, tc.expected, sut.GetV())
		}
	}
}

func TestPPU_IncreaseVByOffset(t *testing.T) {
	sut := NewPPU()
	sut.AppendOnV(0x00)
	sut.AppendOnV(100)

	sut.IncreaseVByOffset(200)
	if sut.GetV() != 300 {
		t.Errorf("expected v to be 300, got %d", sut.GetV())
	}
}

func TestPPU_SetT_GetT(t *testing.T) {
	sut := NewPPU()

	testCases := []struct {
		input    uint16
		expected uint16
	}{
		{input: 0x0000, expected: 0x0000},
		{input: 0x0ABC, expected: 0x0ABC},
		{input: 32767, expected: 0},
		{input: 32768, expected: 1},
	}

	for _, tc := range testCases {
		sut.SetT(tc.input)
		if sut.GetT() != tc.expected {
			t.Errorf("SetT(%d): expected %d, got %d", tc.input, tc.expected, sut.GetT())
		}
	}
}

func TestPPU_SetX_GetX(t *testing.T) {
	sut := NewPPU()

	testCases := []struct {
		input    uint8
		expected uint8
	}{
		{input: 0, expected: 0},
		{input: 3, expected: 3},
		{input: 6, expected: 6},
		{input: 7, expected: 0}, // 7 % 7 == 0
		{input: 8, expected: 1}, // 8 % 7 == 1
	}

	for _, tc := range testCases {
		sut.SetX(tc.input)
		if sut.GetX() != tc.expected {
			t.Errorf("SetX(%d): expected %d, got %d", tc.input, tc.expected, sut.GetX())
		}
	}
}

func TestPPU_FineY(t *testing.T) {
	sut := NewPPU()

	testCases := []struct {
		name     string
		vValue   uint16
		expected uint16
	}{
		{name: "zero v", vValue: 0x0000, expected: 0},
		{name: "fine Y is 1", vValue: 0x1000, expected: 1},
		{name: "fine Y is 2", vValue: 0x2000, expected: 2},
		{name: "fine Y is 7", vValue: 0x7000, expected: 7},
		{name: "fine Y with other bits set", vValue: 0x5ABC, expected: 5},
		{name: "lower bits only", vValue: 0x0FFF, expected: 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut.v = tc.vValue
			fineY := sut.GetFineY()
			if fineY != tc.expected {
				t.Errorf("GetFineY() with v=0x%X: expected %d, got %d", tc.vValue, tc.expected, fineY)
			}
		})
	}
}

func TestPPU_ToggleW_ClearW(t *testing.T) {
	sut := NewPPU()

	if sut.IsWSet() {
		t.Fatal("expected w to start as false")
	}

	sut.ToggleW()
	if !sut.IsWSet() {
		t.Errorf("expected w to be true after first ToggleW")
	}

	sut.ToggleW()
	if sut.IsWSet() {
		t.Errorf("expected w to be false after second ToggleW")
	}

	sut.ToggleW()
	if !sut.IsWSet() {
		t.Errorf("expected w to be true before ClearW")
	}

	sut.ClearW()
	if sut.IsWSet() {
		t.Errorf("expected w to be false after ClearW")
	}

	// Calling ClearW when already false keeps it false
	sut.ClearW()
	if sut.IsWSet() {
		t.Errorf("expected w to remain false on redundant ClearW")
	}
}

func TestPPU_SelectColorByPalleteIndex(t *testing.T) {
	sut := NewPPU()

	indices := []byte{0, 1, 0x0F, 0x10, 0x20, 0x30, 0x3F}
	for _, idx := range indices {
		color := sut.SelectColorByPalleteIndex(idx)
		expected := ColorPallet[idx]
		if color != expected {
			t.Errorf("SelectColorByPalleteIndex(%d): expected 0x%06X, got 0x%06X", idx, expected, color)
		}
	}
}

func TestPPU_ShouldRenderScanlinePixel(t *testing.T) {
	testCases := []struct {
		name     string
		dots     uint16
		scanline uint16
		expected bool
	}{
		{
			name:     "origin point (dots 0, scanline 0) is visible",
			dots:     0,
			scanline: 0,
			expected: true,
		},
		{
			name:     "last visible pixel on visible scanline",
			dots:     MAX_PIXEL_PER_SCANLINE,
			scanline: V_BLANK_SCANLINE_START - 1,
			expected: true,
		},
		{
			name:     "mid visible screen",
			dots:     128,
			scanline: 120,
			expected: true,
		},
		{
			name:     "dots exceed max pixel per scanline",
			dots:     MAX_PIXEL_PER_SCANLINE + 1,
			scanline: 100,
			expected: false,
		},
		{
			name:     "dots at end of scanline",
			dots:     MAX_DOTS_PER_LINE - 1,
			scanline: 50,
			expected: false,
		},
		{
			name:     "visible dots but at VBlank start scanline",
			dots:     0,
			scanline: V_BLANK_SCANLINE_START,
			expected: false,
		},
		{
			name:     "visible dots during VBlank scanlines",
			dots:     100,
			scanline: 250,
			expected: false,
		},
		{
			name:     "pre-render scanline",
			dots:     0,
			scanline: MAX_FRAME_SCANLINE,
			expected: false,
		},
		{
			name:     "both dots and scanline out of visible range",
			dots:     300,
			scanline: 250,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := NewPPU()
			sut.dots = tc.dots
			sut.scanline = tc.scanline

			result := sut.ShouldRenderScanlinePixel()
			if result != tc.expected {
				t.Errorf("ShouldRenderScanlinePixel() with dots=%d, scanline=%d: expected %v, got %v",
					tc.dots, tc.scanline, tc.expected, result)
			}
		})
	}
}

func TestPPU_IsVBlankStarted(t *testing.T) {
	testCases := []struct {
		name     string
		scanline uint16
		pixel    uint8
		expected bool
	}{
		{
			name:     "exact vblank start condition",
			scanline: V_BLANK_SCANLINE_START,
			pixel:    V_BLANK_PIXEL_START,
			expected: true,
		},
		{
			name:     "correct scanline but pixel != 0",
			scanline: V_BLANK_SCANLINE_START,
			pixel:    1,
			expected: false,
		},
		{
			name:     "one scanline before vblank",
			scanline: V_BLANK_SCANLINE_START - 1,
			pixel:    V_BLANK_PIXEL_START,
			expected: false,
		},
		{
			name:     "one scanline after vblank start",
			scanline: V_BLANK_SCANLINE_START + 1,
			pixel:    V_BLANK_PIXEL_START,
			expected: false,
		},
		{
			name:     "frame start",
			scanline: 0,
			pixel:    0,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := NewPPU()
			sut.scanline = tc.scanline
			sut.pixel = tc.pixel

			result := sut.IsVBlankStarted()
			if result != tc.expected {
				t.Errorf("IsVBlankStarted() with scanline=%d, pixel=%d: expected %v, got %v",
					tc.scanline, tc.pixel, tc.expected, result)
			}
		})
	}
}

func TestPPU_IsOnPreRender(t *testing.T) {
	testCases := []struct {
		name     string
		scanline uint16
		pixel    uint8
		expected bool
	}{
		{
			name:     "exact pre-render condition",
			scanline: MAX_FRAME_SCANLINE,
			pixel:    0,
			expected: true,
		},
		{
			name:     "pre-render scanline but pixel != 0",
			scanline: MAX_FRAME_SCANLINE,
			pixel:    1,
			expected: false,
		},
		{
			name:     "one scanline before pre-render",
			scanline: MAX_FRAME_SCANLINE - 1,
			pixel:    0,
			expected: false,
		},
		{
			name:     "first scanline",
			scanline: 0,
			pixel:    0,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := NewPPU()
			sut.scanline = tc.scanline
			sut.pixel = tc.pixel

			result := sut.IsOnPreRender()
			if result != tc.expected {
				t.Errorf("IsOnPreRender() with scanline=%d, pixel=%d: expected %v, got %v",
					tc.scanline, tc.pixel, tc.expected, result)
			}
		})
	}
}

func TestPPU_TogglePPURendering(t *testing.T) {
	testCases := []struct {
		name           string
		value          uint8
		expected       bool
		bgExpected     bool
		spriteExpected bool
	}{
		{
			name:           "disable rendering",
			value:          0b0000_0000,
			expected:       false,
			bgExpected:     false,
			spriteExpected: false,
		},
		{
			name:           "enable rendering",
			value:          0b0001_1000,
			expected:       true,
			bgExpected:     true,
			spriteExpected: true,
		},
		{
			name:           "enable bg rendering; disable sprite rendering",
			value:          0b0000_1000,
			expected:       true,
			bgExpected:     true,
			spriteExpected: false,
		},
		{
			name:           "disable bg rendering; enable sprite rendering",
			value:          0b0001_0000,
			expected:       true,
			bgExpected:     false,
			spriteExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sut := NewPPU()
			sut.RefreshRenderingStateFromLatch(tc.value)

			if sut.IsRenderingEnabled() != tc.expected {
				t.Errorf("IsRenderingEnabled() with value=%08b: expected %v, got %v", tc.value, tc.expected, sut.IsRenderingEnabled())
			}
			if sut.IsBackgroundRenderingEnabled() != tc.bgExpected {
				t.Errorf("IsBackgroundRenderingEnabled() with value=%08b: expected %v, got %v", tc.value, tc.bgExpected, sut.IsBackgroundRenderingEnabled())
			}
			if sut.IsSpriteRenderingEnabled() != tc.spriteExpected {
				t.Errorf("IsSpriteRenderingEnabled() with value=%08b: expected %v, got %v", tc.value, tc.spriteExpected, sut.IsSpriteRenderingEnabled())
			}
		})
	}
}
