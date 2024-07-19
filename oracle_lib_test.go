//go:build cgo && !nolink_libwasmvm

package cosmwasm

import (
	"testing"

	"github.com/ODIN-PROTOCOL/wasmvm/v2/api"
	"github.com/ODIN-PROTOCOL/wasmvm/v2/types"
	"github.com/stretchr/testify/require"
)

const TESTING_ORACLE_MEMORY_LIMIT = 32 * 1024 // KiB

func newTestVM(t *testing.T) (*Vm, func()) {
	vm, err := NewOracleVm(TESTING_ORACLE_MEMORY_LIMIT)
	require.NoError(t, err)

	cleanup := func() {
		api.ReleaseOracleCache(vm.cache)
	}
	return vm, cleanup
}

func TestFailCompileInvalidContent(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	code := []byte("invalid content")
	spanSize := 1 * 1024 * 1024
	_, err := vm.Compile(code, spanSize)
	require.Equal(t, types.ErrValidation, err)
}

func TestRuntimeError(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  i32.const 0
		  i32.const 0
		  i32.div_s
		  drop
		)
		(func)
		(memory 17)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

		`)
	code, _ := vm.Compile(wasm, spanSize)
	_, err := vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrRuntime, err)
}

func TestInvaildSignature(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(func (param i64 i64 i32 i64)
		  (local $idx i32)
		  (local.set $idx (i32.const 0))
		  (block
			  (loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 17)
		(export "prepare" (func 0))
		(export "execute" (func 1)))
	  `)
	code, _ := vm.Compile(wasm, spanSize)
	_, err := vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrBadEntrySignature, err)
}

func TestGasLimit(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (local.set $idx (i32.const 0))
		  (block
			  (loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 17)
		(export "prepare" (func 0))
		(export "execute" (func 1)))
	  `)
	code, err := vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	output, err := vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.NoError(t, err)
	require.Equal(t, api.RunOutput{GasUsed: 70519550000}, output)
	_, err = vm.Prepare(code, 60000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrOutOfGas, err)
}

func TestCompileErrorNoMemory(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (local.set $idx (i32.const 0))
		  (block
			  (loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 10000)))
			  )
			))
		(func)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	code, err := vm.Compile(wasm, spanSize)
	require.Equal(t, types.ErrBadMemorySection, err)
	require.Equal(t, []byte{}, code)
}

func TestCompileErrorMinimumMemoryExceed(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (local.set $idx (i32.const 0))
		  (block
			  (loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 512)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	_, err := vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	wasm = wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (local.set $idx (i32.const 0))
		  (block
			  (loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 513)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	_, err = vm.Compile(wasm, spanSize)
	require.Equal(t, types.ErrBadMemorySection, err)
}

func TestCompileErrorSetMaximumMemory(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(func
		  (local $idx i32)
		  (local.set $idx (i32.const 0))
		  (block
			  (loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 10000)))
			  )
			))
		(func)
		(memory 17 20)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	code, err := vm.Compile(wasm, spanSize)
	require.Equal(t, types.ErrBadMemorySection, err)
	require.Equal(t, []byte{}, code)
}

func TestCompileErrorCheckWasmImports(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(import "env" "beeb" (func (type 0)))
		(func
		(local $idx i32)
		(local.set $idx (i32.const 0))
		(block
				(loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 1000000000)))
				)
			)
		)
		(func)
		(memory 17)
		(data (i32.const 1048576) "beeb")
		(export "prepare" (func 0))
		(export "execute" (func 1)))
		`)
	code, err := vm.Compile(wasm, spanSize)
	require.Equal(t, types.ErrInvalidImports, err)
	require.Equal(t, []byte{}, code)
}

func TestCompileErrorCheckWasmExports(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (func (param i64 i64 i32 i64) (result i64)))
		(import "env" "ask_external_data" (func (type 0)))
		(func
		(local $idx i32)
		(local.set $idx (i32.const 0))
		(block
				(loop
				(local.set $idx (local.get $idx) (i32.const 1) (i32.add) )
				(br_if 0 (i32.lt_u (local.get $idx) (i32.const 1000000000)))
				)
			)
		)
		(memory 17)
		(data (i32.const 1048576) "beeb")
		(export "prepare" (func 0)))
		`)
	code, err := vm.Compile(wasm, spanSize)
	require.Equal(t, types.ErrInvalidExports, err)
	require.Equal(t, []byte{}, code)
}

func TestStackOverflow(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(func call 0)
		(func)
		(memory 10)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	code, _ := vm.Compile(wasm, spanSize)
	_, err := vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrRuntime, err)
}

func TestMemoryGrow(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(func
	i32.const 0
    (memory.grow (i32.const 1))
    i32.gt_s
	if
    	unreachable
    end
     )
		(func)
		(memory 10)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	code, _ := vm.Compile(wasm, spanSize)
	_, err := vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.NoError(t, err)

	wasm = wat2wasm(`(module
		(func
	i32.const 0
    (memory.grow (i32.const 1))
    i32.gt_s
	if
    	unreachable
    end
     )
		(func)
		(memory 512)
		(export "prepare" (func 0))
		(export "execute" (func 1)))

	  `)
	code, _ = vm.Compile(wasm, spanSize)
	_, err = vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrRuntime, err)
}

func TestBadPointer(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (;0;) (func (param i64 i64)))
		(type (;1;) (func))
		(import "env" "set_return_data" (func (;0;) (type 0)))
		(func (type 1)
			i64.const 100000000
			i64.const 1
			call 0
			)
		(func)
		(memory 17)
		(export "prepare" (func 1))
		(export "execute" (func 2)))

		`)
	code, err := vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	_, err = vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrMemoryOutOfBound, err)

	wasm = wat2wasm(`(module
		(type (;0;) (func (param i64 i64 i64 i64)))
		(type (;1;) (func))
		(import "env" "ask_external_data" (func (;0;) (type 0)))
		(func (type 1)
			i64.const 1
			i64.const 1
			i64.const 100000000
			i64.const 1
			call 0
			)
		(func)
		(memory 17)
		(export "prepare" (func 1))
		(export "execute" (func 2)))

		`)
	code, err = vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	_, err = vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrMemoryOutOfBound, err)
}

func TestSpanTooSmall(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (;0;) (func (param i64 i64 i64 i64)))
		(type (;1;) (func))
		(import "env" "ask_external_data" (func (;0;) (type 0)))
		(func (type 1)
			i64.const 1
			i64.const 1
			i64.const 1
			i64.const 1024
			call 0
			)
		(func)
		(memory 17)
		(export "memory" (memory 0))
		(export "prepare" (func 1))
		(export "execute" (func 2)))
		`)
	code, err := vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	_, err = vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.NoError(t, err)

	wasm = wat2wasm(`(module
		(type (;0;) (func (param i64 i64 i64 i64)))
		(type (;1;) (func))
		(import "env" "ask_external_data" (func (;0;) (type 0)))
		(func (type 1)
			i64.const 1
			i64.const 1
			i64.const 1
			i64.const 1025
			call 0
			)
		(func)
		(memory 17)
		(export "memory" (memory 0))
		(export "prepare" (func 1))
		(export "execute" (func 2)))
		`)
	code, err = vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	_, err = vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrSpanTooSmall, err)
}

func TestBadImportSignature(t *testing.T) {
	vm, release := newTestVM(t)
	defer release()

	spanSize := 1 * 1024 * 1024
	wasm := wat2wasm(`(module
		(type (;0;) (func))
		(type (;1;) (func))
		(import "env" "set_return_data" (func (;0;) (type 0)))
		(func
			call 0)
		(func)
		(memory 17)
		(export "memory" (memory 0))
		(export "prepare" (func 1))
		(export "execute" (func 2)))

		`)
	code, err := vm.Compile(wasm, spanSize)
	require.NoError(t, err)
	_, err = vm.Prepare(code, 250000000000, NewMockEnv([]byte(""), 1024))
	require.Equal(t, types.ErrInstantiation, err)
}
