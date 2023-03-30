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

	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp4", Op: OpCheck}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "processors", Value: "cp5", Op: OpCheck}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectSet, Name: "opcodes", Value: "rset", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Name: "opcodes", Op: OpGet}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0/opcodes:rset", T: ObjectSet, Name: "test1", Value: "test2", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0/opcodes:rset", Name: "test1", Op: OpGet}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectMax, Name: "test3", Value: "54", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectMax, Name: "test3", Value: "572", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectMax, Name: "test3", Value: "52", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Name: "test3", Op: OpGet}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectSet, Name: "unaryops", Value: "clr-r0", Op: OpAdd}))
	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectSet, Name: "unaryops", Value: "clr-r2", Op: OpAdd}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", T: ObjectSet, Name: "binaryops", Value: "multf-r2-r0", Op: OpAdd}))

	fmt.Println("----")

	fmt.Println(rg.Requirement(ReqRequest{Node: "/", Op: OpDump}))
	//fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Op: OpDump}))

	//var r ExportedReqs
	//r = make([]ExportedReq, 0)
	//fmt.Println(rg.Export(&r, "/"), r)

	///newrg, _ := Import(&r)
	//fmt.Println(newrg.Requirement(ReqRequest{Node: "/", Op: OpDump}))

	fmt.Println(rg.Requirement(ReqRequest{Node: "/", T: ObjectSet, Name: "prova", Value: "ciao", Op: OpAdd}))
	// fmt.Println(rg.Requirement(ReqRequest{Node: "/processors:cp0", Name: "/prova:ciao", Op: OpClone}))
	rg.Clone("/processors:cp0", "/prova:ciao")
	fmt.Println(rg.Requirement(ReqRequest{Node: "/", Op: OpDump}))

	rg.Close()

}
