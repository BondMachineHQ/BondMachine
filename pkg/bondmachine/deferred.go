package bondmachine

import "errors"

type DeferredInstruction func(*VM) bool

func (vm *VM) AddDeferredInstruction(diName string, di DeferredInstruction) error {
	if vm == nil || vm.DeferredInstructions == nil {
		return errors.New("VM or DeferredInstructions map is nil")
	}
	if _, exists := vm.DeferredInstructions[diName]; !exists {
		vm.DeferredInstructions[diName] = di
	}
	return nil
}

func (vm *VM) ExecuteDeferredInstructions() error {
	if vm == nil || vm.DeferredInstructions == nil {
		return errors.New("VM or DeferredInstructions map is nil")
	}
	notCompleted := make(map[string]DeferredInstruction)
	for diName, di := range vm.DeferredInstructions {
		if complete := di(vm); !complete {
			notCompleted[diName] = di
		}
	}
	vm.DeferredInstructions = notCompleted
	return nil
}
