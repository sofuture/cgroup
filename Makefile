SRC := $(wildcard *.go)
TARGET := main

all: $(TARGET)

$(TARGET): get
	go build -o $@

get: $(SRC)
	go get

clean:
	$(RM) $(TARGET)

test: $(TARGET)
	./main

.PHONY: clean
