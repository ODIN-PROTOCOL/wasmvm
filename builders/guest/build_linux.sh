#!/bin/bash
set -o errexit -o nounset -o pipefail

#yum -y install gmp-dev gmp gmp-devel

build_gnu_x86_64.sh
build_gnu_aarch64.sh
