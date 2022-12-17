package bmreqs

import (
	"fmt"
	"testing"
)

func TestBMReq(t *testing.T) {

	rg := NewReqRoot()

	// Some errors
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Op: 34}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp0", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp1", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp2", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp3", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp4", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", Name: "processors", Op: OpGet}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectSet, Name: "opcodes", Value: "rset", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Name: "opcodes", Op: OpGet}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0/opcodes:rset", T: ObjectSet, Name: "test1", Value: "test2", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0/opcodes:rset", Name: "test1", Op: OpGet}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectMax, Name: "test3", Value: "54", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectMax, Name: "test3", Value: "572", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectMax, Name: "test3", Value: "52", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Name: "test3", Op: OpGet}))

	fmt.Println("----")

	fmt.Println(rg.Requirement(ReqRequest{Node: "/", Op: OpDump}))
	//fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Op: OpDump}))

	rg.Close()

}
