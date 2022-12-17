package bmline

import (
	"fmt"
	"testing"
)

func TestMetaData(t *testing.T) {
	be := new(BasmElement)
	be.SetValue("element")
	be.BasmMeta = be.SetMeta("test3", "value3")
	be.BasmMeta = be.SetMeta("test4", "value4")

	if be.GetMeta("test3") != "value3" {
		t.Errorf("GetMeta failed")
	}
	if be.GetMeta("test4") != "value4" {
		t.Errorf("GetMeta failed")
	}

	bl := new(BasmLine)
	bl.Operation = be
	bl.Elements = make([]*BasmElement, 0)
	bl.BasmMeta = bl.SetMeta("test2", "value2")

	if bl.GetMeta("test2") != "value2" {
		t.Errorf("GetMeta failed")
	}

	bb := new(BasmBody)
	bb.Lines = make([]*BasmLine, 1)
	bb.BasmMeta = bb.SetMeta("test1", "value1")
	bb.Lines[0] = bl

	if bb.GetMeta("test1") != "value1" {
		t.Errorf("GetMeta failed")
	}

	pref := new(BasmElement)
	pref.BasmMeta = pref.SetMeta("test2", "prefix")

	fmt.Println(bb)
	bb.PrefixMeta(pref)
	fmt.Println(bb)

	// fmt.Println("bl test2", bl.GetMeta("test2"))
	// fmt.Println("be test3", be.GetMeta("test3"))

	// bb.RmMeta("test1")
	// bl.RmMeta("test1")
	// be.RmMeta("test1")

	// fmt.Println("bb test1", bb.GetMeta("test1"))
	// fmt.Println("bl test2", bl.GetMeta("test2"))
	// fmt.Println("be test3", be.GetMeta("test3"))

	// bbc := bb.Copy()
	// bll := bl.Copy()
	// bec := be.Copy()

	// fmt.Printf("bbc pointer %p\n", bbc)
	// fmt.Printf("bll pointer %p\n", bll)
	// fmt.Printf("bec pointer %p\n", bec)

	// fmt.Println("bbc test1", bbc.GetMeta("test1"))
	// fmt.Println("bll test2", bll.GetMeta("test2"))
	// fmt.Println("bec test3", bec.GetMeta("test3"))

	// fmt.Println("be List", be.ListMeta())
}
