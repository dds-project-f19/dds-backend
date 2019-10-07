if not exist build mkdir build
cd build

set GOOS = "windows"
set GOARCH = "amd64"
go build -o server_x64_windows.exe ..

set GOOS = "windows"
set GOARCH = "386"
go build -o server_x86_windows.exe ..

set GOOS = "linux"
set GOARCH = "amd64"
go build -o server_x64_linux ..

set GOOS = "linux"
set GOARCH = "386"
go build -o server_x86_linux ..

xcopy "..\config.yaml"