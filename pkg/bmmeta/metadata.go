package bmmeta

import (
	"sort"
	"strings"
)

// Metadata
type BasmMeta struct {
	metaData map[string]string
}

func (bm *BasmMeta) LoopMeta() map[string]string {
	if bm != nil && bm.metaData != nil {
		return bm.metaData
	}
	return map[string]string{}
}

func (bm *BasmMeta) ListMeta() string {
	result := ""
	if bm != nil && bm.metaData != nil {
		// Order the keys
		keys := make([]string, 0, len(bm.metaData))
		for k := range bm.metaData {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			val := bm.metaData[key]
			result += "," + key + ":" + val
		}
	}
	if len(result) != 0 {
		return result[1:]
	}
	return result
}

func (bm *BasmMeta) GetMeta(meta string) string {
	if bm != nil && bm.metaData != nil {
		if value, ok := bm.metaData[meta]; ok {
			return value
		} else {
			return ""
		}
	}
	return ""
}

func (bm *BasmMeta) RmMeta(meta string) {
	if bm != nil && bm.metaData != nil {
		if _, ok := bm.metaData[meta]; ok {
			delete(bm.metaData, meta)
		}
	}
}

func (bm *BasmMeta) SetMeta(meta string, value string) *BasmMeta {
	if bm == nil {
		newbm := new(BasmMeta)
		newbm.metaData = make(map[string]string)
		newbm.metaData[meta] = value
		return newbm
	}
	bm.metaData[meta] = value
	return bm

}

func (bm *BasmMeta) AddMeta(meta string, value string) *BasmMeta {
	if bm == nil {
		newbm := new(BasmMeta)
		newbm.metaData = make(map[string]string)
		newbm.metaData[meta] = value
		return newbm
	}
	if _, ok := bm.metaData[meta]; ok {
		splitted := strings.Split(bm.metaData[meta], ":")
		if !stringInSlice(value, splitted) {
			splitted = append(splitted, value)
			sort.Strings(splitted)
			bm.metaData[meta] = strings.Join(splitted, ":")
		}
	} else {
		bm.metaData[meta] = value
	}
	return bm

}
