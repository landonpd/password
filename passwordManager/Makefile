TARGET=password
FILE=main.go

all: build 


build:
	go build -o bin/$(TARGET) $(FILE)

buildWin:
# echo "Compiling for every OS and Platform"
# GOOS=darwin GOARCH=386 go build -o bin/windows/$(TARGET) main.go mac theoretically
# GOOS=linux GOARCH=386 go build -o bin/linux/$(TARGET) main.go

	GOOS=windows GOARCH=386 go build -o bin/$(TARGET).exe $(FILE)

buildAll:
	GOOS=windows GOARCH=386 go build -o bin/$(TARGET).exe $(FILE)
	go build -o bin/$(TARGET) $(FILE)


run:
	go run $(FILE)

clean:
	rm -r bin
	go build -o bin/$(TARGET) $(FILE)