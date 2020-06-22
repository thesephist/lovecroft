lovecroft = ./src/lovecroft.go

all: run


run:
	go run -race ${lovecroft}


# build for specific OS target
build-%:
	GOOS=$* GOARCH=amd64 go build -o lovecroft-$* ${lovecroft}


build:
	go build -o lovecroft ${lovecroft}


# clean any generated files
clean:
	rm -rvf lovecroft lovecroft-*
