export GOOS=windows
export GOARCH=amd64
cd ..
echo " Building ..."
go build -o build/windows/slacktui.exe
echo "  Build done: build/windows/slacktui.exe"
