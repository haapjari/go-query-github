include make.properties

run: 
	go run ${MAIN_MODULE}

build:
	go build -o bin/main ${MAIN_MODULE}