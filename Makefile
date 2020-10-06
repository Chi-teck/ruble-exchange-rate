TARGET  = rer
SOURCES = $(sort $(wildcard src/*.c))
OBJECTS = $(SOURCES:.c=.o)
HEADERS = $(sort $(wildcard src/*.h))
OPTIONS = -std=gnu99 -O2 -Wall -lcurl -lxml2 -I/usr/include/libxml2

${TARGET}: ${OBJECTS}
	gcc -o ${TARGET} ${OBJECTS} ${OPTIONS}

src/main.o: src/main.c ${HEADERS}
	gcc -c src/main.c ${OPTIONS} -o src/main.o

src/options.o: src/options.c ${HEADERS}
	gcc -c src/options.c -o src/options.o

src/fetch.o: src/fetch.c ${HEADERS}
	gcc -c src/fetch.c ${OPTIONS} -o src/fetch.o

src/decoder.o: src/decoder.c ${HEADERS}
	gcc -c src/decoder.c ${OPTIONS} -o src/decoder.o
	
install:
	install $(TARGET) /usr/local/bin/$(TARGET)
	
remove:
	rm /usr/local/bin/$(TARGET)
	
clean:
	rm -f ${OBJECTS} $(TARGET)

.PHONY: clean install remove
