SRC := $(wildcard *.go)
TARGET := main

all: $(TARGET)

$(TARGET): $(SRC)
	go build -o $@

clean:
	$(RM) $(TARGET)

test: $(TARGET)
	./main

.PHONY: clean
