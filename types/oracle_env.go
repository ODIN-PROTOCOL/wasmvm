package types

type EnvInterface interface {
	GetSpanSize() int64
	GetCalldata() []byte
	SetReturnData([]byte) error
	GetAskCount() int64
	GetMinCount() int64
	GetPrepareTime() int64
	GetExecuteTime() (int64, error)
	GetAnsCount() (int64, error)
	AskExternalData(eid int64, did int64, data []byte) error
	GetExternalDataStatus(eid int64, vid int64) (int64, error)
	GetExternalData(eid int64, vid int64) ([]byte, error)
}
