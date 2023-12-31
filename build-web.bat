set GOOS=js
set GOARCH=wasm
go build -o red-cat-run-2d.wasm -ldflags="-s -w" .
copy "C:\Program Files\Go\misc\wasm\wasm_exec.js" red-cat-run-2d.js
set GOOS=
set GOARCH=

del red-cat-run-2d.zip
7z a -tzip red-cat-run-2d.zip index.html red-cat-run-2d.js red-cat-run-2d.wasm
copy index.html ..\hinst.github.io\red-cat-run-2d\index.html
copy red-cat-run-2d.js ..\hinst.github.io\red-cat-run-2d\red-cat-run-2d.js
copy red-cat-run-2d.wasm ..\hinst.github.io\red-cat-run-2d\red-cat-run-2d.wasm

del red-cat-run-2d.js
del red-cat-run-2d.wasm
