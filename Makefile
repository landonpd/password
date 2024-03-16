TARGET=password
FILE=main.go

ifeq  ($(shell echo "Windows"), "Windows")
	TARGET := $(TARGET).exe
	DEL = del /q /s /f
else
	DEL = rm -rf
endif

all: 
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
	$(DEL) bin
	go build -o bin/$(TARGET) $(FILE)