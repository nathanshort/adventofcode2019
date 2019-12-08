package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"os"
)

const (
	black       = '0'
	white       = '1'
	transparent = '2'
	width       = 25
	height      = 6
	size        = width * height
)

type Layer struct {
	/// map of pixel value to count of pixels with that value
	ValueCounts map[byte]int
	Data        []byte
	Width       int
	Height      int
}

func (l Layer) save(writer io.Writer) {
	image := image.NewRGBA(image.Rect(0, 0, l.Width, l.Height))
	for h := 0; h < l.Height; h++ {
		for w := 0; w < l.Width; w++ {
			pixel := w + h*l.Width
			switch l.Data[pixel] {
			case black:
				image.Set(w, h, color.RGBA{0, 0, 0, 255})
			case white:
				image.Set(w, h, color.RGBA{255, 255, 255, 255})
			}
		}
	}
	png.Encode(writer, image)
}

func newEmptyLayer(width int, height int) *Layer {
	l := new(Layer)
	l.Width = width
	l.Height = height
	l.ValueCounts = make(map[byte]int)
	l.Data = make([]byte, width*height)
	return l
}

func newLayer(data []byte, width int, height int) *Layer {
	l := newEmptyLayer(width, height)
	l.ValueCounts = make(map[byte]int)
	for _, color := range data {
		l.ValueCounts[color] = l.ValueCounts[color] + 1
	}
	l.Data = data
	return l
}

func part1(layers []*Layer) {
	minZeros := math.MaxInt32
	valueAtMinZeros := 0
	for l := 0; l < len(layers); l++ {
		if layers[l].ValueCounts['0'] < minZeros {
			minZeros = layers[l].ValueCounts['0']
			valueAtMinZeros = layers[l].ValueCounts['1'] * layers[l].ValueCounts['2']
		}
	}
	fmt.Printf("part 1: %d\n", valueAtMinZeros)
}

func part2(layers []*Layer) {
	composite := newEmptyLayer(width, height)
	for l := 0; l < len(layers); l++ {
		for h := 0; h < height; h++ {
			for w := 0; w < width; w++ {
				pixel := w + h*width
				if l == 0 || composite.Data[pixel] == transparent {
					composite.Data[pixel] = layers[l].Data[pixel]
				}
			}
		}
	}
	f, _ := os.Create("/tmp/image")
	defer f.Close()
	composite.save(f)
}

func main() {
	input, _ := ioutil.ReadAll(os.Stdin)
	var layers []*Layer
	for start := 0; start < len(input); start += size {
		layer := newLayer(input[start:start+size], width, height)
		layers = append(layers, layer)
	}

	part1(layers)
	part2(layers)

}
