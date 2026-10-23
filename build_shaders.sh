#!/usr/bin/env bash
#
# Builds the SPIR-V shaders (Vulkan) from shader_src/*.hlsl into shaders/.
#
# DXIL (Direct3D 12) is built on Windows with build_shaders.bat: DXIL has to be
# signed by dxil.dll, which only exists there, so this script doesn't make it.
#
# dxc needs SPIR-V support. The easiest way to get one is the Vulkan SDK from
# LunarG (run its setup-env.sh so VULKAN_SDK and PATH are set). You can also
# point DXC at any dxc binary:  DXC=/path/to/dxc ./build_shaders.sh
#
set -euo pipefail
cd "$(dirname "$0")"

DXC="${DXC:-}"
if [ -z "$DXC" ]; then
    if command -v dxc >/dev/null 2>&1; then
        DXC=dxc
    elif [ -n "${VULKAN_SDK:-}" ] && [ -x "$VULKAN_SDK/bin/dxc" ]; then
        DXC="$VULKAN_SDK/bin/dxc"
    else
        echo "dxc not found. Install the Vulkan SDK (and source its setup-env.sh), or set DXC=/path/to/dxc." >&2
        exit 1
    fi
fi
echo "using $DXC"

SRC=shader_src
OUT=shaders
mkdir -p "$OUT"

build() {
    local name=$1
    "$DXC" -spirv -fspv-target-env=vulkan1.0 -T vs_6_0 -E main -Fo "$OUT/$name.vert.spv" "$SRC/$name.vert.hlsl"
    "$DXC" -spirv -fspv-target-env=vulkan1.0 -T ps_6_0 -E main -Fo "$OUT/$name.frag.spv" "$SRC/$name.frag.hlsl"
}

for name in lit shadow tonemap; do
    build "$name"
done

# Optional: validate the output if spirv-val is installed (package spirv-tools).
if command -v spirv-val >/dev/null 2>&1; then
    for file in "$OUT"/*.spv; do
        if ! spirv-val --target-env vulkan1.0 "$file"; then
            echo "validation failed: $file" >&2
        fi
    done
fi

echo done
