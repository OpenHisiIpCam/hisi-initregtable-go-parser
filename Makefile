all: clean build

systemrdl:
	make -C ./registers-description build

generate: systemrdl
	./rdl_to_go.py ./registers-description/hi3516av200.rdl hi3516av200 > ./hi3516av200.go
	go fmt ./hi3516av200.go
	./rdl_to_go.py ./registers-description/hi3516av200.rdl hi3519v101 > ./hi3519v101.go
	go fmt ./hi3519v101.go
	./rdl_to_go.py ./registers-description/hi3520dv200.rdl hi3520dv200 > ./hi3520dv200.go
	go fmt ./hi3520dv200.go

build: generate
	go build

clean:
	rm -f ./hi35*.go
	rm -f ./hisi-initregtable-go-parser

deps:
	pip3 install systemrdl-compiler
