use crate::env::{Env, RunOutput};

use crate::span::Span;
use crate::vm::VMQuerier;

use failure::{bail, Error as FailureError};
use std::panic::catch_unwind;

use owasm_vm::cache::{Cache, CacheOptions};
use owasm_vm::error::Error;

// Cache initializing section
#[repr(C)]
pub struct oracle_cache_t {}

pub fn to_cache(ptr: *mut oracle_cache_t) -> Option<&'static mut Cache> {
    if ptr.is_null() {
        None
    } else {
        let c = unsafe { &mut *(ptr as *mut Cache) };
        Some(c)
    }
}

#[no_mangle]
pub extern "C" fn oracle_init_cache(size: u32) -> *mut oracle_cache_t {
    let r = catch_unwind(|| oracle_do_init_cache(size)).unwrap_or_else(|_| bail!("Caught panic"));
    match r {
        Ok(t) => t as *mut oracle_cache_t,
        Err(_) => std::ptr::null_mut(),
    }
}

fn oracle_do_init_cache(size: u32) -> Result<*mut Cache, FailureError> {
    let cache = Cache::new(CacheOptions { cache_size: size });
    let out = Box::new(cache);
    Ok(Box::into_raw(out))
}

#[no_mangle]
pub unsafe extern "C" fn oracle_release_cache(cache: *mut oracle_cache_t) {
    if !cache.is_null() {
        // this will free cache when it goes out of scope
        let _ = Box::from_raw(cache as *mut Cache);
    }
}

// Compile and execute section
#[no_mangle]
pub extern "C" fn do_compile(input: Span, output: &mut Span) -> Error {
    match owasm_vm::compile(input.read()) {
        Ok(out) => {
            output.write(&out);
            Error::NoError
        }
        Err(e) => e,
    }
}

#[no_mangle]
pub extern "C" fn do_run(
    cache: *mut oracle_cache_t,
    code: Span,
    gas_limit: u64,
    is_prepare: bool,
    env: Env,
    output: &mut RunOutput,
) -> Error {
    if !cache.is_null() {
        let vm_querier = VMQuerier::new(env);
        match owasm_vm::run(to_cache(cache).unwrap(), code.read(), gas_limit, is_prepare, vm_querier) {
            Ok(gas_used) => {
                output.gas_used = gas_used;
                Error::NoError
            }
            Err(e) => e,
        }
    } else {
        Error::UnknownError
    }
}
