#!/usr/bin/bash

# Helper script to install Go dev tools.
# This is run _inside_ the manylinux container(s)
# used in cibuildwheel to build the wheels.

set -euo pipefail

TDIR="$(mktemp -d)"
>&2 echo "Working dir: ${TDIR}"
trap "rm -rf ${TDIR}" EXIT

>&2 echo "Downloading Go 1.23.6 distribution file."
curl -fL -o "${TDIR}/go1.23.6.linux-amd64.tar.gz" 'https://go.dev/dl/go1.23.6.linux-amd64.tar.gz'

>&2 echo "Checking distribution file integrity"
GO_DIST_SHA256='9379441ea310de000f33a4dc767bd966e72ab2826270e038e78b2c53c2e7802d'
printf '%s %s/go1.23.6.linux-amd64.tar.gz\n' "${GO_DIST_SHA256}" "${TDIR}" | sha256sum -c

>&2 echo "Unpacking to /usr/local/go"
rm -rf /usr/local/go && tar -C /usr/local -xzf "${TDIR}/go1.23.6.linux-amd64.tar.gz"

>&2 echo "Installed Go version:"
/usr/local/go/bin/go version
