package bminfo

import (
	"encoding/json"
	"io/ioutil"
)

type BMinfo struct {
	List map[string]string
}

// WriteBMinfo writes a BMinfo struct to a file JSON encoded.

func (bmInfo *BMinfo) WriteBMinfo(filename string) error {
	b, err := json.Marshal(bmInfo)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadBMinfo reads a BMinfo struct from a file JSON encoded
func ReadBMinfo(filename string) (*BMinfo, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	bmInfo := new(BMinfo)
	err = json.Unmarshal(b, bmInfo)
	if err != nil {
		return nil, err
	}
	return bmInfo, nil
}
