#ifndef CWRAPPER_H
#define CWRAPPER_H


typedef struct {
    double xval, yval;
}cBasicPoint;

typedef struct {
	cBasicPoint p1,p2;
}cEPosition;

typedef struct {
	cEPosition freePosition;
	float centerAngle;
	int centerDistance;
}cEmptyPosRelation;

typedef struct {
		cEPosition parent;
		int eastSize,southSize,westSize,northSize,southeastSize,southwestSize,northeastSize,northwestSize;
        cEmptyPosRelation *east,*south,*west,*north,*southeast,*southwest,*northeast,*northwest;
}cEmptyPosRelations;

typedef struct {
	int count;
	cEmptyPosRelations *exports;
	unsigned char *modimg; 
	int imgsize;
}cExports;

 #endif