package cosmwasm

import "C"
import "github.com/CosmWasm/wasmvm/v2/internal/api"

type Vm struct {
	cache api.OracleCache
}

func NewOracleVm(size uint32) (*Vm, error) {
	cache, err := api.InitOracleCache(size)

	if err != nil {
		return nil, err
	}

	return &Vm{
		cache: cache,
	}, nil
}

func (vm Vm) Compile(code []byte, spanSize int) ([]byte, error) {
	return api.Compile(code, spanSize)
}

func (vm Vm) Prepare(code []byte, gasLimit uint64, env api.EnvInterface) (api.RunOutput, error) {
	return api.OracleRun(vm.cache, code, gasLimit, true, env)
}

func (vm Vm) Execute(code []byte, gasLimit uint64, env api.EnvInterface) (api.RunOutput, error) {
	return api.OracleRun(vm.cache, code, gasLimit, false, env)
}
