package bmreqs

import (
	"errors"
	"fmt"
	"strconv"
)

type objectMax struct {
	name  string
	t     uint8
	value int64
}

func (o *objectMax) init() {
}

//

func (o *objectMax) getName() string {
	return o.name
}

func (o *objectMax) getType() uint8 {
	return o.t
}

func (o *objectMax) setName(name string) error {
	o.name = name
	return nil
}

func (o *objectMax) setType(t uint8) error {
	o.t = t
	return nil
}

//

func (o *objectMax) insertReq(req string) error {
	if i, err := strconv.ParseInt(req, 10, 64); err == nil {
		if i > o.value {
			o.value = i
		}
		return nil
	}
	return errors.New("Integer parse failed")
}

func (o *objectMax) removeReq(req string) error {
	return errors.New("Remove request not implemented in objectMax")
}

//

func (o *objectMax) getReqs() string {
	return fmt.Sprint(o.value)
}

//

func (o *objectMax) supportSub() bool {
	return false
}

func (o *objectMax) listSub() []string {
	return []string{}
}

func (o *objectMax) getSub(req string) (*bmReqObj, error) {
	return nil, nil
}
