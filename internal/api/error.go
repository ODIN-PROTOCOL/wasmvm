package api

// #include "bindings.h"
import "C"
import (
	"github.com/ODIN-PROTOCOL/wasmvm/v2/types"
)

// toCError converts the given Golang error into Rust/C error.
func toCError(err error) C.Error {
	switch err {
	case nil:
		return C.Error_NoError
	case types.ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	case types.ErrTooManyExternalData:
		return C.Error_TooManyExternalDataError
	case types.ErrDuplicateExternalID:
		return C.Error_DuplicateExternalIDError
	case types.ErrBadValidatorIndex:
		return C.Error_BadValidatorIndexError
	case types.ErrBadExternalID:
		return C.Error_BadExternalIDError
	case types.ErrUnavailableExternalData:
		return C.Error_UnavailableExternalDataError
	case types.ErrRepeatSetReturnData:
		return C.Error_RepeatSetReturnDataError
	default:
		return C.Error_UnknownError
	}
}

// toGoError converts the given Rust/C error to Golang error.
func toGoError(code C.Error) error {
	switch code {
	case C.Error_NoError:
		return nil
	case C.Error_SpanTooSmallError:
		return types.ErrSpanTooSmall
	// Rust-generated errors during compilation.
	case C.Error_ValidationError:
		return types.ErrValidation
	case C.Error_DeserializationError:
		return types.ErrDeserialization
	case C.Error_SerializationError:
		return types.ErrSerialization
	case C.Error_InvalidImportsError:
		return types.ErrInvalidImports
	case C.Error_InvalidExportsError:
		return types.ErrInvalidExports
	case C.Error_BadMemorySectionError:
		return types.ErrBadMemorySection
	case C.Error_GasCounterInjectionError:
		return types.ErrGasCounterInjection
	case C.Error_StackHeightInjectionError:
		return types.ErrStackHeightInjection
	// Rust-generated errors during runtime.
	case C.Error_InstantiationError:
		return types.ErrInstantiation
	case C.Error_RuntimeError:
		return types.ErrRuntime
	case C.Error_OutOfGasError:
		return types.ErrOutOfGas
	case C.Error_BadEntrySignatureError:
		return types.ErrBadEntrySignature
	case C.Error_MemoryOutOfBoundError:
		return types.ErrMemoryOutOfBound
	// Go-generated errors while interacting with OEI.
	case C.Error_WrongPeriodActionError:
		return types.ErrWrongPeriodAction
	case C.Error_TooManyExternalDataError:
		return types.ErrTooManyExternalData
	case C.Error_DuplicateExternalIDError:
		return types.ErrDuplicateExternalID
	case C.Error_BadValidatorIndexError:
		return types.ErrBadValidatorIndex
	case C.Error_BadExternalIDError:
		return types.ErrBadExternalID
	case C.Error_UnavailableExternalDataError:
		return types.ErrUnavailableExternalData
	case C.Error_RepeatSetReturnDataError:
		return types.ErrRepeatSetReturnData
	default:
		return types.ErrUnknown
	}
}
