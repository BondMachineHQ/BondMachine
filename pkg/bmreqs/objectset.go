package bmreqs

import (
	"errors"
	"fmt"
	"strings"
)

type objectSet struct {
	name string
	t    uint8
	set  map[string]*bmReqObj
}

func (o *objectSet) init() {
	o.set = make(map[string]*bmReqObj)
}

//

func (o *objectSet) getName() string {
	return o.name
}

func (o *objectSet) getType() uint8 {
	return o.t
}

func (o *objectSet) setName(name string) error {
	o.name = name
	return nil
}

func (o *objectSet) setType(t uint8) error {
	o.t = t
	return nil
}

//

func (o *objectSet) insertReq(req string) error {
	if o.set == nil {
		return fmt.Errorf("uninitialized Set")
	}
	if _, ok := o.set[req]; !ok {
		newObj := new(bmReqObj)
		newObj.init()
		o.set[req] = newObj
	}
	return nil
}

func (o *objectSet) removeReq(req string) error {
	if o.set != nil {
		if _, ok := o.set[req]; ok {
			delete(o.set, req)
		}
	} else {
		return errors.New("uninitialized Set")
	}
	return nil
}

//

func (o *objectSet) checkReq(req string) (string, error) {
	if o.set != nil {
		if _, ok := o.set[req]; ok {
			return "true", nil
		} else {
			return "false", nil
		}
	}
	return "", errors.New("uninitialized Set")
}

//

func (o *objectSet) getReqs() string {
	if o.set == nil {
		return ""
	}
	keys := make([]string, 0, len(o.set))
	for k := range o.set {
		keys = append(keys, k)
	}
	return fmt.Sprint(strings.Join(keys, ","))
}

func (o *objectSet) importReqs(rg *ReqRoot, node string, name string, req string) error {
	for _, r := range strings.Split(req, ",") {
		resp := rg.Requirement(ReqRequest{Node: node, T: ObjectSet, Name: name, Value: r, Op: OpAdd})
		if resp.Error != nil {
			return resp.Error
		}
	}
	return nil
}

//

func (o *objectSet) supportSub() bool {
	return true
}

func (o *objectSet) listSub() []string {
	if o.set == nil {
		return []string{}
	}

	keys := make([]string, 0, len(o.set))
	for k := range o.set {
		keys = append(keys, k)
	}
	return keys
}

func (o *objectSet) getSub(req string) (*bmReqObj, error) {
	if o.set == nil {
		return nil, errors.New("uninitialized Set")
	}
	if node, ok := o.set[req]; ok {
		return node, nil
	}
	return nil, nil
}
