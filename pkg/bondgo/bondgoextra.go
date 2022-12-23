package bondgo

import (
	"fmt"
	"sync"
)

type inputs []*Input

var link_inputs map[int]inputs
var link_outputs map[int]int
var curr_in int
var curr_out int

var wg sync.WaitGroup

type Input struct {
	idx  int
	iidx int
	val  interface{}
}

type Output struct {
	idx  int
	oidx int
}

func (o *Output) Write(val interface{}) {
	for idx, oidx := range link_outputs {
		if oidx == o.oidx {
			for _, ii := range link_inputs[idx] {
				ii.val = val
			}
		}
	}
}

func (o *Output) Make(idx int) {
	o.oidx = curr_out
	o.idx = idx
	curr_out++
	if _, ok := link_outputs[idx]; ok {
		panic("Output already defined")
	}
	link_outputs[idx] = o.oidx
	//show_internal()
	wg.Done()
}

func (i *Input) Read() interface{} {
	return i.val
}

func (i *Input) Make(idx int, val interface{}) {
	i.iidx = curr_in
	i.idx = idx
	i.val = val
	curr_in++
	if linp, ok := link_inputs[idx]; ok {
		link_inputs[idx] = append(linp, i)
	} else {
		link_inputs[idx] = make([]*Input, 1)
		link_inputs[idx][0] = i
	}
	//show_internal()
	wg.Done()
}

func AllWait() {
	wg.Wait()
}

func AllInit(gornum int) {
	wg.Add(gornum)
}

func Void(interface{}) {
}

func show_internal() {
	fmt.Println("link_inputs", link_inputs)
	fmt.Println("link_outputs", link_outputs)
	fmt.Println("curr_in", curr_in)
	fmt.Println("curr_out", curr_out)
}

func init() {
	link_outputs = make(map[int]int)
	link_inputs = make(map[int]inputs)
	//show_internal()
}
