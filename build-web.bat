set GOOS=js
set GOARCH=wasm
go build -o red-cat-run-2d.wasm .
copy "C:\Program Files\Go\misc\wasm\wasm_exec.js" red-cat-run-2d.js
set GOOS=
set GOARCH=