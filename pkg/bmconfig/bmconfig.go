package bmconfig

const (
	DisableDynamicalMatching = uint64(0) + iota
	ChooserMinWordSize
)

type BmConfig struct {
	config map[uint64]struct{}
}

func NewBmConfig() *BmConfig {
	bmc := new(BmConfig)
	bmc.config = make(map[uint64]struct{})
	return bmc
}

func (bmc *BmConfig) Activate(key uint64) {
	bmc.config[key] = struct{}{}
}

func (bmc *BmConfig) Deactivate(key uint64) {
	if _, ok := bmc.config[key]; !ok {
		return
	}
	delete(bmc.config, key)
}

func (bmc *BmConfig) IsActive(key uint64) bool {
	if _, ok := bmc.config[key]; !ok {
		return false
	}
	return true
}
