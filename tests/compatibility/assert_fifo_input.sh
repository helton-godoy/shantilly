#!/bin/bash
#
# GUI widgets for shell scripts - SHantilly version 1.0
#
# Copyright (C) 2015-2016, 2020 Andriy Martynets <andy.martynets@gmail.com>
# This file is part of SHantilly and is distributed under the terms of the
# GNU General Public License version 3 or later.

set -euo pipefail

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 EXECUTABLE EXPECTED_FILE" >&2
    exit 2
fi

executable=$1
expected_file=$2
test_dir=$(mktemp -d)
fifo_path=${test_dir}/commands.fifo
output_path=${test_dir}/stdout.txt
error_path=${test_dir}/stderr.txt
shantilly_pid=

cleanup() {
    exec 3>&- || true
    if [[ -n ${shantilly_pid} ]]; then
        kill "${shantilly_pid}" 2>/dev/null || true
        wait "${shantilly_pid}" 2>/dev/null || true
    fi
    rm -rf "${test_dir}"
}
trap cleanup EXIT

mkfifo "${fifo_path}"
QT_QPA_PLATFORM=offscreen "${executable}" --hidden <"${fifo_path}" >"${output_path}" 2>"${error_path}" &
shantilly_pid=$!

exec 3>"${fifo_path}"
printf '%s\n' 'add checkbox "FIFO" fifo_value checked' >&3
printf '%s\n' 'query' >&3

for _ in {1..40}; do
    if grep -Fxq 'fifo_value=1' "${output_path}"; then
        break
    fi
    sleep 0.05
done

if ! kill -0 "${shantilly_pid}" 2>/dev/null; then
    echo "SHantilly exited before the second FIFO write" >&2
    sed 's/^/  /' "${error_path}" >&2
    exit 1
fi
printf '%s\n' 'unset fifo_value checked' >&3
printf '%s\n' 'query' >&3

for _ in {1..40}; do
    if grep -Fxq 'fifo_value=0' "${output_path}"; then
        break
    fi
    sleep 0.05
done

if ! cmp -s "${expected_file}" "${output_path}"; then
    echo "Unexpected FIFO output" >&2
    diff -u "${expected_file}" "${output_path}" >&2 || true
    echo "SHantilly stderr:" >&2
    sed 's/^/  /' "${error_path}" >&2
    exit 1
fi
