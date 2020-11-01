#include <opencv2/core/core.hpp>
#include <opencv2/highgui/highgui.hpp>
#include <iostream>
#include "mapProvider.hpp"
#include "cWrapper.h"
#include <stdlib.h>
#include <stdio.h>

using namespace cv;
using namespace std;


extern int testCExport(cExports exp);


int main( int argc, char** argv )
{
    if( argc != 2)
    {
     cout <<" Usage: display_image ImageToLoadAndDisplay" << endl;
     return -1;
    }


/*
    Mat image;
    image = imread(argv[1], CV_LOAD_IMAGE_COLOR);   // Read the file
    imshow("asda",image);
    waitKey(0);
    if(! image.data )                              // Check for invalid input
    {
        cout <<  "Could not open or find the image" << std::endl ;
        return -1;
    }
    int rows = image.rows;
    int cols = image.cols;
*/

    FILE *f = fopen(argv[1], "rb");
    fseek(f, 0, SEEK_END);
    long fsize = ftell(f);
    fseek(f, 0, SEEK_SET);  //same as rewind(f);
    unsigned char *string = (unsigned char*)malloc(fsize + 1);
    fread(string, fsize, 1, f);
    fclose(f);

    cExports toRet = basic(string,fsize);

    //testCExport(toRet);
    //std::vector<char> data(toRet.modimg, toRet.modimg + toRet.imgsize);
    //Mat image = imdecode(Mat(data), 1);
    //imshow ("asda",image);

    waitKey(0);                                          // Wait for a keystroke in the window
    return 0;
}



