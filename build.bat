setlocal
pushd %~dp0

set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w"

popd

