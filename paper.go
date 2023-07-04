package paper

import (
	"image"
	"image/color"
)

func New(th Theme, w, h int) *Paper {
	return &Paper{
		theme: th,
		bound: image.Rect(0, 0, w, h),
	}
}

type Theme func(byte) color.Color

func mix(prim, bg, v byte) byte {
	return byte((int(255-v)*int(prim) + int(v)*int(bg)) / 255)
}

func defineTheme(prim, bg color.RGBA) Theme {
	return func(v byte) color.Color {
		return color.RGBA{
			R: mix(prim.R, bg.R, v),
			G: mix(prim.G, bg.G, v),
			B: mix(prim.B, bg.B, v),
			A: 255,
		}
	}
}

var (
	Modern Theme = defineTheme(
		color.RGBA{21, 21, 21, 255},
		color.RGBA{221, 221, 221, 255},
	)
	Nostalgia Theme = defineTheme(
		color.RGBA{45, 40, 14, 255},
		color.RGBA{227, 218, 189, 255},
	)
	Sepia Theme = defineTheme(
		color.RGBA{52, 36, 36, 255},
		color.RGBA{190, 155, 118, 255},
	)
	Night Theme = defineTheme(
		color.RGBA{221, 221, 221, 255},
		color.RGBA{21, 21, 21, 255},
	)
)

type Paper struct {
	theme  Theme
	bound  image.Rectangle
	tiles  []*tile
	masked bool
}

func (p *Paper) ColorModel() color.Model {
	if p.masked {
		return color.GrayModel
	}
	return color.RGBAModel
}

func (p *Paper) Bounds() image.Rectangle {
	return p.bound
}

func (p *Paper) ExtendWidth(extra int) {
	p.bound.Max.X += extra
}

func (p *Paper) ExtendHeight(extra int) {
	p.bound.Max.Y += extra
}

func (p *Paper) Mask() {
	p.masked = true
}

func (p *Paper) Unmask() {
	p.masked = false
}

func (p *Paper) At(x, y int) color.Color {
	v := byte(255)
	if image.Pt(x, y).In(p.bound) {
		i := 0
		for i < len(p.tiles) && y >= p.tiles[i].bound.Max.Y {
			i++
		}
		if i < len(p.tiles) && y >= p.tiles[i].bound.Min.Y {
			for t := p.tiles[i]; t != nil; t = t.next {
				if x < t.bound.Min.X {
					break
				}
				if x < t.bound.Max.X {
					v = t.Get(x, y)
					break
				}
			}
		}
	}
	if p.masked {
		return color.Gray{v}
	}
	return p.theme(v)
}

func (p *Paper) Set(x, y int, c color.Color) {
	if !image.Pt(x, y).In(p.bound) {
		return
	}

	r, g, b, _ := c.RGBA()
	v := byte((19595*r + 38470*g + 7471*b + 1<<15) >> 24)

	bx, by := (x/8)*8, (y/8)*8
	bound := image.Rect(bx, by, bx+8, by+8)

	i := 0
	for i < len(p.tiles) && y >= p.tiles[i].bound.Max.Y {
		i++
	}

	t := &tile{bound: bound}
	if i == len(p.tiles) {
		p.tiles = append(p.tiles, t)
	} else if y < p.tiles[i].bound.Min.Y {
		p.tiles = append(p.tiles, t)
		copy(p.tiles[i+1:], p.tiles[i:])
		p.tiles[i] = t
	} else {
		t = p.link(t, i, x)
	}

	t.Set(x, y, v)
}

func (p *Paper) link(t *tile, i, x int) *tile {
	var (
		u *tile
		v = p.tiles[i]
	)
	for v != nil && x >= v.bound.Max.X {
		u, v = v, v.next
	}

	if v == nil {
		u.next = t
	} else if x < v.bound.Max.X && x >= v.bound.Min.X {
		t = v
	} else if u == nil {
		t.next = v
		p.tiles[i] = t
	} else {
		t.next = v
		u.next = t
	}

	return t
}

type tile struct {
	bound  image.Rectangle
	pixels [64]byte
	next   *tile
}

func (t *tile) Get(x, y int) byte {
	xo, yo := x-t.bound.Min.X, y-t.bound.Min.Y
	return 255 - t.pixels[8*yo+xo]
}

func (t *tile) Set(x, y int, v byte) {
	xo, yo := x-t.bound.Min.X, y-t.bound.Min.Y
	t.pixels[8*yo+xo] = 255 - v
}
