#!/usr/bin/env bash

set -euo pipefail

version="${1:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
dist_dir="${2:-dist}"

targets=(
  darwin/amd64
  darwin/arm64
  linux/amd64
  linux/arm64
  windows/amd64
  windows/arm64
)

rm -rf "${dist_dir}"
mkdir -p "${dist_dir}"

for target in "${targets[@]}"; do
  os="${target%/*}"
  arch="${target#*/}"
  artifact="notel_${version}_${os}_${arch}"
  binary_name="notel"

  if [[ "${os}" == "windows" ]]; then
    binary_name="notel.exe"
  fi

  mkdir -p "${dist_dir}/${artifact}"
  GOOS="${os}" GOARCH="${arch}" go build \
    -ldflags "-X github.com/stennie/notel/cmd.Version=${version}" \
    -o "${dist_dir}/${artifact}/${binary_name}" .

  if [[ "${os}" == "windows" ]]; then
    (
      cd "${dist_dir}"
      zip -qr "${artifact}.zip" "${artifact}"
    )
  else
    tar -C "${dist_dir}" -czf "${dist_dir}/${artifact}.tar.gz" "${artifact}"
  fi
done

checksum_file="${dist_dir}/SHA256SUMS"

if command -v shasum >/dev/null 2>&1; then
  checksum_cmd=(shasum -a 256)
elif command -v sha256sum >/dev/null 2>&1; then
  checksum_cmd=(sha256sum)
else
  echo "no SHA-256 checksum tool found (expected shasum or sha256sum)" >&2
  exit 1
fi

(
  cd "${dist_dir}"
  shopt -s nullglob
  for artifact in *.tar.gz *.zip; do
    "${checksum_cmd[@]}" "${artifact}"
  done
) > "${checksum_file}"
