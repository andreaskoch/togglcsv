install:
	go run make.go -install

crosscompile:
	go run make.go -crosscompile

test:
	go run make.go -test

coverage:
	go run make.go -coverage
