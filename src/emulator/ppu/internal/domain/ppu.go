package domain

const (
	PPU_CONTROL = 0x2000
	PPU_MASK    = 0x2001
	PPU_STATUS  = 0x2002
	OAM_ADDRESS = 0x2003
	OAM_DATA    = 0x2004
	PPU_SCROLL  = 0x2005
	PPU_ADDRESS = 0x2006
	PPU_DATA    = 0x2007
	OAM_DMA     = 0x4014
)

const (
	MAX_FRAME_SCANLINE     = 261
	MAX_PIXEL_PER_SCANLINE = 255
	V_BLANK_SCANLINE_START = 240
	V_BLANK_SCANLINE_END   = 260
	V_BLANK_PIXEL_START    = 0
	MAX_DOTS_PER_LINE      = 336
)

const BASE_NAMETABLE_ADDRESS = 0x2000

const (
	MAX_VALUE_FOR_15_BITS = 32767
	MAX_VALUE_FOR_3_BITS  = 7
)

var ColorPallet = [64]uint32{
	0x7C7C7C, 0x0000FC, 0x0000BC, 0x4428BC, 0x940084, 0xA80020, 0xA81000, 0x881400,
	0x503000, 0x007800, 0x006800, 0x005800, 0x004058, 0x000000, 0x000000, 0x000000,
	0xBCBCBC, 0x0070EC, 0x3C40FC, 0x7C00FA, 0xA800B4, 0xC43800, 0xE04000, 0xE86814,
	0x9C9400, 0x44B400, 0x34C400, 0x00D840, 0x00BCBC, 0x000000, 0x000000, 0x000000,
	0x3CBCFC, 0x0078F8, 0x0058F8, 0x6844FC, 0xD800CC, 0xE40058, 0xF83800, 0xE45C10,
	0xAC7C00, 0x00B800, 0x00A800, 0x00A844, 0x008888, 0x000000, 0x000000, 0x000000,
	0xF8F8F8, 0xA4E4FC, 0xB8B8F8, 0xD8B8F8, 0xF8B8F8, 0xF8A4C0, 0xF0D0B0, 0xFCE0A8,
	0xE8D878, 0xD8F878, 0xB8F8B8, 0xB8F8D8, 0x00FCFC, 0xD8D8D8, 0x000000, 0x000000,
}

type PPU struct {
	scanline     uint16
	pixel        uint8
	enableRender bool
	v, t         uint16
	x            uint8
	w            bool
	dots         uint16
}

func NewPPU() *PPU {
	p := &PPU{
		enableRender: false,
	}

	return p
}

func (p *PPU) AdvanceToNextScanlinePixel() {
	p.pixel++
}

func (p *PPU) AdvanceToNextScanline() {
	if p.dots == MAX_DOTS_PER_LINE-1 {
		p.scanline++
		p.scanline %= (MAX_FRAME_SCANLINE + 1)
	}
}

func (p *PPU) IncreaseVByOffset(offset uint16) {
	p.setV(p.v + offset)
}

func (p *PPU) IncreaseDots() {
	p.dots = (p.dots + 1) % MAX_DOTS_PER_LINE
}

func (p *PPU) SelectColorByPalleteIndex(index byte) uint32 {
	return ColorPallet[index]
}

func (p *PPU) ShouldRenderScanlinePixel() bool {
	return p.dots <= MAX_PIXEL_PER_SCANLINE && p.scanline < V_BLANK_SCANLINE_START
}

func (p *PPU) IsVBlankStarted() bool {
	return p.scanline == V_BLANK_SCANLINE_START && p.pixel == V_BLANK_PIXEL_START
}

func (p *PPU) IsOnPreRender() bool {
	return p.scanline == MAX_FRAME_SCANLINE && p.pixel == 0
}

func (p *PPU) GetCurrentDots() uint16 {
	return p.dots
}

func (p *PPU) GetCurrentScanline() uint16 {
	return p.scanline
}

func (p *PPU) GetCurrentScanlinePixel() uint8 {
	return p.pixel
}

func (p *PPU) GetFineY() uint16 {
	const fine_y_mask = 0x7000
	return (p.v & fine_y_mask) >> 12
}

func (p *PPU) AppendOnV(value byte) {
	const high_byte_mask = 0x3F00

	if !p.IsWSet() {
		p.SetT((uint16(value) << 8) & high_byte_mask)
	} else {
		p.SetT(p.GetT() | uint16(value))
		p.copyTToV()
	}

	p.ToggleW()
}

func (p *PPU) GetV() uint16 {
	return p.v
}

func (p *PPU) SetT(value uint16) {
	p.t = value % MAX_VALUE_FOR_15_BITS
}

func (p *PPU) GetT() uint16 {
	return p.t
}

func (p *PPU) SetX(value uint8) {
	p.x = value % MAX_VALUE_FOR_3_BITS
}

func (p *PPU) GetX() uint8 {
	return p.x
}

func (p *PPU) ToggleW() {
	p.w = !p.w
}

func (p *PPU) ClearW() {
	p.w = false
}

func (p *PPU) IsWSet() bool {
	return p.w
}

func (p *PPU) IsRenderingEnabled() bool {
	return p.enableRender
}

func (p *PPU) setV(value uint16) {
	p.v = value % MAX_VALUE_FOR_15_BITS
}

func (p *PPU) copyTToV() {
	p.setV(p.t)
}
