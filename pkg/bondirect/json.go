package bondirect

import (
	"encoding/json"
	"os"
)

func UnmarshalMesh(c *Config, meshFile string) (*Mesh, error) {

	mesh := new(Mesh)

	if _, err := os.Stat(meshFile); err == nil {
		if jsonFile, err := os.ReadFile(meshFile); err == nil {
			if err := json.Unmarshal([]byte(jsonFile), mesh); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}

	return mesh, nil
}
