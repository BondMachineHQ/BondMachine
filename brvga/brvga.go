package brvga

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type cpTextMemory struct {
	cpId    int
	width   int
	height  int
	topPos  int
	leftPos int
	mem     []byte
}

type BrvgaTextMemory struct {
	constraintString string
	cps              []cpTextMemory
}

func NewBrvgaTextMemory(constraint string) (*BrvgaTextMemory, error) {
	components := strings.Split(constraint, ":")
	componentsN := len(components)
	if (componentsN-1)%5 == 0 {
		boxes := make([]cpTextMemory, 0)

		for i := 1; i < componentsN; i = i + 5 {
			newBox := cpTextMemory{}
			if newCP, err := strconv.Atoi(components[i]); err == nil {
				newBox.cpId = newCP
			} else {
				return nil, err
			}
			if newLeft, err := strconv.Atoi(components[i+1]); err == nil {
				newBox.leftPos = newLeft
			} else {
				return nil, err
			}
			if newTop, err := strconv.Atoi(components[i+2]); err == nil {
				newBox.topPos = newTop
			} else {
				return nil, err
			}
			if newWidth, err := strconv.Atoi(components[i+3]); err == nil {
				newBox.width = newWidth
			} else {
				return nil, err
			}
			if newHeight, err := strconv.Atoi(components[i+4]); err == nil {
				newBox.height = newHeight
			} else {
				return nil, err
			}

			memL := newBox.width * newBox.height
			newBox.mem = make([]byte, memL)
			for j := 0; j < memL; j++ {
				newBox.mem[j] = 0x00
			}

			boxes = append(boxes, newBox)
		}

		result := new(BrvgaTextMemory)
		result.constraintString = constraint
		result.cps = boxes
		return result, nil
	}

	return nil, errors.New("invalid constraint string")
}

func (b *BrvgaTextMemory) Dump() string {
	result := ""
	result += fmt.Sprintf("Constraint %s\n\n", b.constraintString)
	for _, cp := range b.cps {
		result += fmt.Sprintf("cp %d: %d x %d at %d, %d\n", cp.cpId, cp.width, cp.height, cp.leftPos, cp.topPos)
		for i := 0; i < cp.height*cp.width; i++ {
			if i%cp.width == 0 {
				result += "\n"
			}
			result += fmt.Sprintf("%02x ", cp.mem[i])
		}
		result += "\n\n"
	}
	return result
}