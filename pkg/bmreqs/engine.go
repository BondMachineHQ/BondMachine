package bmreqs

import (
	"errors"
	"fmt"
	"strings"
)

func (rg *ReqRoot) decodeNode(node string) (*bmReqObj, error) {
	if node == "/" {
		return rg.obj, nil
	}

	seqlist := strings.Split(node, "/")
	var currNode *bmReqObj = rg.obj
	var lName string
	var lReq string
	for _, leaf := range seqlist {
		currMap := currNode.getMap()
		if leaf != "" {
			decLeaf := strings.Split(leaf, ":")
			if len(decLeaf) != 2 {
				return nil, errors.New("malformed node")
			}
			lReq = decLeaf[0]
			lName = decLeaf[1]

			if reqset, ok := currMap[lReq]; ok {
				if reqset.supportSub() {
					if newNode, err := reqset.getSub(lName); err == nil {
						if newNode != nil {
							currNode = newNode
						} else {
							return nil, errors.New("node empty")
						}
					} else {
						return nil, errors.New("node not exists")
					}
				} else {
					return nil, errors.New("requirement Set does not support sub requirements")
				}
			} else {
				return nil, errors.New("unknown Requirement Set")
			}
		}
	}

	return currNode, nil
}

func (rg *ReqRoot) run() {
	ctx := rg.ctx
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-rg.ask:
			resp := ReqResponse{Value: "", Error: nil}
			switch req.Op {
			case OpAdd:
				if node, err := rg.decodeNode(req.Node); err == nil {
					currMap := node.getMap()
					if reqSet, ok := currMap[req.Name]; ok {
						if reqSet.getType() != req.T {
							resp.Error = errors.New("insert failed: Mismatch types")
						} else {
							if err := reqSet.insertReq(req.Value); err != nil {
								resp.Error = errors.New("Insert failed: " + fmt.Sprint(err))
							}
						}
					} else {
						switch req.T {
						case ObjectSet:
							newReq := new(objectSet)
							newReq.init()
							newReq.setName(req.Name)
							newReq.setType(ObjectSet)
							if err := newReq.insertReq(req.Value); err == nil {
								node.getMap()[req.Name] = newReq
							} else {
								resp.Error = errors.New("insert failed: " + fmt.Sprint(err))
							}
						case ObjectMax:
							newReq := new(objectMax)
							newReq.init()
							newReq.setName(req.Name)
							newReq.setType(ObjectMax)
							if err := newReq.insertReq(req.Value); err == nil {
								node.getMap()[req.Name] = newReq
							} else {
								resp.Error = errors.New("insert failed: " + fmt.Sprint(err))
							}
						default:
							resp.Error = errors.New("unknown Type")
						}
					}
				} else {
					resp.Error = errors.New("node decoding failed: " + fmt.Sprint(err))
				}
			case OpCheck:
				if node, err := rg.decodeNode(req.Node); err == nil {
					currMap := node.getMap()
					if reqSet, ok := currMap[req.Name]; ok {
						if reqSet.getType() != req.T {
							resp.Error = errors.New("check failed: Mismatch types")
						} else {
							if present, err := reqSet.checkReq(req.Value); err != nil {
								resp.Error = errors.New("check failed: " + fmt.Sprint(err))
							} else {
								resp.Value = present
							}
						}
					} else {
					}
				} else {
					resp.Error = errors.New("node decoding failed: " + fmt.Sprint(err))
				}
			case OpGet:
				if node, err := rg.decodeNode(req.Node); err == nil {
					if reqSet, ok := node.getMap()[req.Name]; ok {
						resp.Value = reqSet.getReqs()
					} else {
						resp.Error = errors.New("set of requirements not found \"" + req.Name + "\"")
					}
				} else {
					resp.Error = errors.New("Node decoding failed: " + fmt.Sprint(err))
				}
			case OpDump:
				if result, err := rg.recursiveDump(req.Node); err == nil {
					resp.Value = result
				} else {
					resp.Error = err
				}
			case OpClone:
				go rg.Clone(req.Node, req.Name)
				resp.Value = "Clone started"
			default:
				resp.Error = errors.New("unknown Operation")
			}

			rg.answer <- resp
		}
	}
}

func (rg *ReqRoot) recursiveDump(node string) (string, error) {
	if n, err := rg.decodeNode(node); err == nil {
		result := ""
		for name, set := range n.bmReqMap {
			result += node[1:] + "/" + name + "[" + set.getReqs() + "]\n"
			if set.supportSub() {
				subs := set.listSub()
				for _, sub := range subs {
					sr, _ := rg.recursiveDump(node + "/" + name + ":" + sub)
					list := strings.Split(sr, "\n")
					for _, l := range list {
						if l != "" {
							result += l + "\n"
						}
					}
				}
			}
		}
		return result, nil
	} else {
		return "", errors.New("node decoding failed: " + fmt.Sprint(err))
	}
}

func (rg *ReqRoot) Export(r *ExportedReqs, node string) error {
	if n, err := rg.decodeNode(node); err == nil {
		for name, set := range n.bmReqMap {
			*r = append(*r, ExportedReq{Node: node, Type: set.getType(), Req: set.getReqs(), Name: name})
			if set.supportSub() {
				subs := set.listSub()
				for _, sub := range subs {
					if node == "/" {
						node = node[1:]
					}
					rg.Export(r, node+"/"+name+":"+sub)
				}
			}
		}
		return nil
	} else {
		return errors.New("node decoding failed: " + fmt.Sprint(err))
	}
}

func Import(r *ExportedReqs) (*ReqRoot, error) {
	rg := NewReqRoot()
	for _, req := range *r {
		switch req.Type {
		case ObjectSet:
			n := new(objectSet)
			if err := n.importReqs(rg, req.Node, req.Name, req.Req); err != nil {
				return nil, err
			}
		case ObjectMax:
			n := new(objectMax)
			if err := n.importReqs(rg, req.Node, req.Name, req.Req); err != nil {
				return nil, err
			}
		}
	}
	return rg, nil
}

func (rg *ReqRoot) Clone(node string, nodeDes string) error {
	if n, err := rg.decodeNode(node); err == nil {
		for name, set := range n.bmReqMap {
			if node == "/" {
				node = node[1:]
			}
			switch set.getType() {
			case ObjectSet:
				n := new(objectSet)
				if err := n.importReqs(rg, nodeDes, name, set.getReqs()); err != nil {
					return err
				}
			case ObjectMax:
				n := new(objectMax)
				if err := n.importReqs(rg, nodeDes, name, set.getReqs()); err != nil {
					return err
				}
			}
			if set.supportSub() {
				subs := set.listSub()
				for _, sub := range subs {
					if err := rg.Clone(node+"/"+name+":"+sub, nodeDes+"/"+name+":"+sub); err != nil {
						return err
					}
				}
			}
		}
		return nil
	} else {
		return errors.New("node decoding failed: " + fmt.Sprint(err))
	}
}
