//go:build cgo && !nolink_libwasmvm

package cosmwasm

import (
	"github.com/ODIN-PROTOCOL/wasmvm/v2/internal/api"
)

func libwasmvmVersionImpl() (string, error) {
	return api.LibwasmvmVersion()
}
