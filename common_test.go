package cosmwasm

import (
	wasmtime "github.com/bytecodealliance/wasmtime-go/v20"
)

// wat2wasm compiles the given Wat content to Wasm, relying on the host's wat2wasm program.
func wat2wasm(wat string) []byte {
	wasm, err := wasmtime.Wat2Wasm(wat)
	if err != nil {
		panic(err)
	}

	return wasm
}

type RawRequest struct {
	ExternalID   int64
	DataSourceID int64
	Calldata     []byte
}

func NewRawRequest(eid int64, did int64, calldata []byte) RawRequest {
	return RawRequest{
		ExternalID:   eid,
		DataSourceID: did,
		Calldata:     calldata,
	}
}

type RawReport struct {
	ExternalID int64
	ExitCode   uint32
	Data       []byte
}

type OracleMockEnv struct {
	Calldata    []byte
	Retdata     []byte
	rawRequests []RawRequest
	SpanSize    int64
}

func NewMockEnv(calldata []byte, spanSize int64) *OracleMockEnv {
	return &OracleMockEnv{
		Calldata:    calldata,
		Retdata:     []byte{},
		rawRequests: []RawRequest{},
		SpanSize:    spanSize,
	}
}

func (env *OracleMockEnv) GetSpanSize() int64 {
	return env.SpanSize
}

func (env *OracleMockEnv) GetCalldata() []byte {
	return env.Calldata
}

func (env *OracleMockEnv) SetReturnData(data []byte) error {
	env.Retdata = data
	return nil
}

func (env *OracleMockEnv) AskExternalData(eid int64, did int64, data []byte) error {
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		eid, did, data,
	))
	return nil
}

func (env *OracleMockEnv) GetExternalDataFull(eid int64, valIdx int64) ([]byte, int64) {
	return []byte("BEEB"), 0
}

func (env *OracleMockEnv) GetExternalDataStatus(eid int64, vid int64) (int64, error) {
	_, status := env.GetExternalDataFull(eid, vid)
	return status, nil
}

func (env *OracleMockEnv) GetExternalData(eid int64, vid int64) ([]byte, error) {
	data, _ := env.GetExternalDataFull(eid, vid)
	return data, nil
}

func (env *OracleMockEnv) GetAskCount() int64 {
	return 0
}

func (env *OracleMockEnv) GetMinCount() int64 {
	return 0
}

func (env *OracleMockEnv) GetPrepareTime() int64 {
	return 0
}

func (env *OracleMockEnv) GetExecuteTime() (int64, error) {
	return 0, nil
}

func (env *OracleMockEnv) GetAnsCount() (int64, error) {
	return 0, nil
}
