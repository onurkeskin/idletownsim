CC = `g++`
CFLAGS = `pkg-config --cflags opencv`
LIBS = `pkg-config --libs opencv`
LDLIBS = `pkg-config --libs opencv`
DEPS = mapProvider.hpp BasicPoint.hpp e_position.hpp utils.hpp EmptyPosRelation.hpp cWrapper.h
OBJ = test.o mapProvider.o e_position.o BasicPoint.o utils.o EmptyPosRelation.o
PROGRAMS = test


%.o: %.c $(DEPS)
	$(CC) -I/usr/local/include/opencv -I/usr/local/include/opencv2 -L/usr/local/lib/ $(CFLAGS) $(LIBS) -c -fPIC -o $@ $<

test: $(OBJ)
	g++ -o $@ $^ $(CFLAGS) $(LIBS)

lib: lib
	g++ -I/usr/local/include/opencv -I/usr/local/include/opencv2 -L/usr/local/lib/ $(CFLAGS) $(LIBS) -shared -o libmapProvider.so -fPIC mapProvider.cpp e_position.cpp BasicPoint.cpp utils.cpp EmptyPosRelation.cpp

clean: 
	rm -f $(PROGRAMS) *.o *.so core