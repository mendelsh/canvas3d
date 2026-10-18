@echo off
setlocal
cd /d "%~dp0"

set SRC=shader_src
set OUT=shaders
if not exist %OUT% mkdir %OUT%

dxc -T vs_6_0 -E main -Fo %OUT%\lit.vert.dxil %SRC%\lit.vert.hlsl || exit /b 1
dxc -T ps_6_0 -E main -Fo %OUT%\lit.frag.dxil %SRC%\lit.frag.hlsl || exit /b 1
dxc -T vs_6_0 -E main -Fo %OUT%\shadow.vert.dxil %SRC%\shadow.vert.hlsl || exit /b 1
dxc -T ps_6_0 -E main -Fo %OUT%\shadow.frag.dxil %SRC%\shadow.frag.hlsl || exit /b 1

@REM  dxc -spirv -fspv-target-env=vulkan1.0 -T vs_6_0 -E main -Fo %OUT%\lit.vert.spv %SRC%\lit.vert.hlsl || exit /b 1
@REM  dxc -spirv -fspv-target-env=vulkan1.0 -T ps_6_0 -E main -Fo %OUT%\lit.frag.spv %SRC%\lit.frag.hlsl || exit /b 1

echo done
