all:
	@echo "TODO"

systemrdl:
	make -C ./registers-description build

generate: systemrdl
	rdl_to_go.py ./registers-description/hi3516av200.rdl

build:
	go build

deps:
	pip3 install systemrdl-compiler
