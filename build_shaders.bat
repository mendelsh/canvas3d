@echo off
setlocal
cd /d "%~dp0"

rem Find dxc: PATH first, then the Vulkan SDK, then the Windows SDK.
set DXC=
where dxc >nul 2>nul && set DXC=dxc
if not defined DXC if defined VULKAN_SDK if exist "%VULKAN_SDK%\Bin\dxc.exe" set DXC="%VULKAN_SDK%\Bin\dxc.exe"
set "WINKITS=%ProgramFiles(x86)%\Windows Kits\10\bin"
if not defined DXC for /f "delims=" %%P in ('where /r "%WINKITS%" dxc.exe 2^>nul') do echo %%P | find "\x64\" >nul && set DXC="%%P"
if not defined DXC (
    echo dxc not found. Install the Windows SDK or the Vulkan SDK, or put dxc.exe on PATH.
    exit /b 1
)
echo using %DXC%

set SRC=shader_src
set OUT=shaders
if not exist %OUT% mkdir %OUT%

for %%S in (lit shadow tonemap) do call :build %%S || exit /b 1

echo done
exit /b 0

:build
%DXC% -T vs_6_0 -E main -Fo %OUT%\%1.vert.dxil %SRC%\%1.vert.hlsl || exit /b 1
%DXC% -T ps_6_0 -E main -Fo %OUT%\%1.frag.dxil %SRC%\%1.frag.hlsl || exit /b 1

rem SPIR-V is optional: only dxc builds with SPIR-V support (the Vulkan SDK one) can make it.
%DXC% -spirv -fspv-target-env=vulkan1.0 -T vs_6_0 -E main -Fo %OUT%\%1.vert.spv %SRC%\%1.vert.hlsl >nul 2>nul || echo skipped SPIR-V for %1.vert
%DXC% -spirv -fspv-target-env=vulkan1.0 -T ps_6_0 -E main -Fo %OUT%\%1.frag.spv %SRC%\%1.frag.hlsl >nul 2>nul || echo skipped SPIR-V for %1.frag
exit /b 0
