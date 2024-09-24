include ../../Makefile.Common

build:
	make -C wasm build-wasm
	make -C web bundle-js


