package brvgasdl

import (
	"bufio"
	"encoding/hex"
	"os"

	"github.com/BondMachineHQ/BondMachine/brvga"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
)

type Header struct {
	header []byte
}

type Fonts struct {
	fonts  []byte
	images map[byte]canvas.Image
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

	// Read the fonts file
	fontsFile, err := os.Open(fontsPath)
	if err != nil {
		return nil, err
	}
	defer fontsFile.Close()

	scanner := bufio.NewScanner(fontsFile)
	for scanner.Scan() {
		if line := scanner.Text(); len(line) > 0 {

			if line[:2] == "0x" {
				hexString := line[2:]
				if decoded, err := hex.DecodeString(hexString); err == nil {
					result.Fonts.fonts = append(result.Fonts.fonts, decoded...)
				}
			}
		}
	}

	// Load 8x8 fonts
	for c := 0; c < 128; c++ {
		for i := 0; i < 8; i++ {
			// charLine := result.Fonts.fonts[c*8+i]
		}
	}

	return result, nil
}

func (b *BrvgaSdl) Run() {
	wnd := b.Window
	cv := b.Canvas
	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#000000")
		cv.FillRect(0, 0, w, h)

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
