package bmbuilder

import (
	"errors"
)

const (
	passMetaExtractor = uint64(1) << iota
	passGeneratorsExec
	LAST_PASS = passGeneratorsExec
)

func getPassFunction() map[uint64]func(*BMBuilder) error {
	return map[uint64]func(*BMBuilder) error{
		passMetaExtractor:  metaExtractor,
		passGeneratorsExec: generatorsExec,
	}
}

func getPassFunctionName() map[uint64]string {
	return map[uint64]string{
		passMetaExtractor:  "metaExtractor",
		passGeneratorsExec: "generatorsExec",
	}
}

func IsOptionalPass() map[uint64]bool {
	return map[uint64]bool{
		passMetaExtractor:  false,
		passGeneratorsExec: true,
	}
}

func (bi *BMBuilder) ActivePass(active uint64) bool {
	return (bi.passes & active) != uint64(0)
}

func GetPassMnemonic() map[uint64]string {
	return map[uint64]string{
		passMetaExtractor:  "metaextractor",
		passGeneratorsExec: "generatorsexec",
	}

}

func (bi *BMBuilder) SetActive(pass string) error {
	for passN, v := range GetPassMnemonic() {
		if v == pass {
			if ch, ok := IsOptionalPass()[passN]; ok {
				if ch {
					bi.passes = bi.passes | passN
					return nil
				} else {
					return errors.New("pass is not optional")
				}
			} else {
				return errors.New("pass is not defined")
			}
		}
	}
	return errors.New("pass not found")
}

func (bi *BMBuilder) UnsetActive(pass string) error {
	for passN, v := range GetPassMnemonic() {
		if v == pass {
			if ch, ok := IsOptionalPass()[passN]; ok {
				if ch {
					bi.passes = bi.passes & ^passN
					return nil
				} else {
					return errors.New("pass is not optional")
				}
			} else {
				return errors.New("pass is not defined")
			}
		}
	}
	return errors.New("pass not found")
}
