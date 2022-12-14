package brvgasdl

import (
	"bufio"
	"context"
	"encoding/hex"
	"image"
	"image/color"
	"os"

	"github.com/BondMachineHQ/BondMachine/pkg/brvga"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

type Header struct {
	header []byte
}

type Fonts struct {
	fonts  []byte
	images map[byte]*image.RGBA
}

type BrvgaSdl struct {
	*brvga.BrvgaTextMemory
	*sdlcanvas.Window
	*canvas.Canvas
	*Header
	*Fonts
}

func NewBrvgaSdlUnixSock(constraint string, sockPath string, headerPath string, fontsPath string) (*BrvgaSdl, error) {
	result := new(BrvgaSdl)

	// Create the brvga memory
	textMem, err := brvga.NewBrvgaTextMemory(constraint)
	if err != nil {
		return nil, err
	}
	result.BrvgaTextMemory = textMem

	// Create the canvas and window
	wnd, canvas, err := sdlcanvas.CreateWindow(800, 600, "BondMachine")
	if err != nil {
		return nil, err
	}
	result.Window = wnd
	result.Canvas = canvas

	// Allocate the fonts array
	result.Fonts = new(Fonts)
	result.Fonts.fonts = make([]byte, 0)

	// Allocate the header array
	result.Header = new(Header)
	result.Header.header = make([]byte, 0)

	// Read the header
	headerFile, err := os.Open(headerPath)
	if err != nil {
		return nil, err
	}
	defer headerFile.Close()

	headerReader := bufio.NewScanner(headerFile)
	for headerReader.Scan() {
		if line := headerReader.Text(); len(line) > 0 {
			if line[:2] == "0x" {
				hexString := line[2:]
				if decoded, err := hex.DecodeString(hexString); err == nil {
					result.header = append(result.header, decoded...)
				}
			}
		}
	}

	// Read the fonts file
	fontsFile, err := os.Open(fontsPath)
	if err != nil {
		return nil, err
	}
	defer fontsFile.Close()

	fontsReader := bufio.NewScanner(fontsFile)
	for fontsReader.Scan() {
		if line := fontsReader.Text(); len(line) > 0 {

			if line[:2] == "0x" {
				hexString := line[2:]
				if decoded, err := hex.DecodeString(hexString); err == nil {
					result.Fonts.fonts = append(result.Fonts.fonts, decoded...)
				}
			}
		}
	}

	result.images = make(map[byte]*image.RGBA)

	// Load 8x8 fonts
	for c := 0; c < 128; c++ {

		newChar := image.NewRGBA(image.Rectangle{Max: image.Point{X: 8, Y: 8}})

		for i := 0; i < 8; i++ {
			charLine := result.Fonts.fonts[c*8+i]

			for j := 0; j < 8; j++ {
				if charLine&byte((1<<j)) > 0 {
					newChar.Set(j, i, color.White)
				} else {
					// newChar.Set(j, i, color.Black)
				}
			}
		}

		result.Fonts.images[byte(c)] = newChar
	}

	ctx, _ := context.WithCancel(context.Background())

	go result.UNIXSockReceiver(ctx, "/tmp/brvga.sock")

	return result, nil
}

func (b *BrvgaSdl) Run() {
	wnd := b.Window
	cv := b.Canvas
	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#111")
		cv.FillRect(0, 0, w, h)
		cv.SetFillStyle("#000000")
		cv.FillRect(0, 0, 800, 600)

		for i := 0; i < len(b.header); i++ {
			cv.DrawImage(b.Fonts.images[b.header[i]], 1+float64(i*8), 1)
		}

		offsetX := 6
		offsetY := 14

		for i, cp := range b.BrvgaTextMemory.Cps {
			cv.SetFillStyle("#333")
			cv.FillRect(8*float64(cp.LeftPos)+float64(offsetX), 8*float64(cp.TopPos)+float64(offsetY), 8*float64(cp.Width), 8*float64(cp.Height))
			switch i % 4 {
			case 0:
				cv.SetStrokeStyle("#01c")
			case 1:
				cv.SetStrokeStyle("#075")
			case 2:
				cv.SetStrokeStyle("#a3d")
			case 3:
				cv.SetStrokeStyle("#603")
			}
			cv.StrokeRect(8*float64(cp.LeftPos)-2+float64(offsetX), 8*float64(cp.TopPos)-2+float64(offsetY), 8*float64(cp.Width)+4, 8*float64(cp.Height)+4)
		}

		cv.SetStrokeStyle("#ffffff")
		for _, cp := range b.BrvgaTextMemory.Cps {
			if mem, err := b.GetCpMem(cp.CpId); err == nil {
				posX := 8*float64(cp.LeftPos) + float64(offsetX)
				posY := 8*float64(cp.TopPos) + float64(offsetY)
				for i := 0; i < len(mem); i++ {
					cv.DrawImage(b.Fonts.images[mem[i]], posX+float64(i%cp.Width)*8, posY+float64(i/cp.Width)*8)
				}
			}
		}

		// for r := 0.0; r < math.Pi*2; r += math.Pi * 0.1 {
		// 	cv.SetFillStyle(int(r*10), int(r*20), int(r*40))
		// 	cv.BeginPath()
		// 	cv.MoveTo(w*0.5, h*0.5)
		// 	cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, r, r+0.1*math.Pi, false)
		// 	cv.ClosePath()
		// 	cv.Fill()
		// }

		// cv.SetStrokeStyle("#FFF")
		// cv.SetLineWidth(10)
		// cv.BeginPath()
		// cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.4, 0, math.Pi*2, false)
		cv.Stroke()
	})

}

func (b *BrvgaSdl) Close() {
	b.Window.Destroy()
}
