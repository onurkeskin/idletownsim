/*
void parseMap(int, void*){

    Mat templ(40, 40, CV_8UC3, Scalar(0,0,0));

 	Mat img_display;
  	croppedImg.copyTo( img_display );

	int result_cols =  croppedImg.cols - templ.cols + 1;
 	int result_rows = croppedImg.rows - templ.rows + 1;

	Mat result;
  	result.create( result_rows, result_cols, CV_32FC1 );
 	matchTemplate( croppedImg, templ, result, match_method );
 	//normalize( result, result, 0, 1, NORM_MINMAX, -1, Mat() );

	double minVal; double maxVal; Point minLoc; Point maxLoc;
  	Point matchLoc;

  	printf("Printing for k = spaceCount:%d\n",spaceCount);

  	Scalar grayScalar = Scalar::all(128);
  	Vec3b grayPix = Vec3b(128,128,128);

  	for(int k=1;k <= spaceCount;k++)
  	{
    	minMaxLoc( result, &minVal, &maxVal, &minLoc, &maxLoc, Mat() );

		//printf("minLoc: %d,%d | maxLoc: %d,%d \n",*sLocX,*sLocY,*fLocX,*fLocY);

    	//result.at<float>(minLoc.x,minLoc.y)=1.0;
    	//result.at<float>(maxLoc.x,maxLoc.y)=0.0;

  		/// For SQDIFF and SQDIFF_NORMED, the best matches are lower values. For all the other methods, the higher the better
  		if( match_method  == CV_TM_SQDIFF || match_method == CV_TM_SQDIFF_NORMED )
    		{ matchLoc = minLoc; }
  			else
    		{ matchLoc = maxLoc; }

    	int* LocX = &(matchLoc.x);
    	int* LocY = &(matchLoc.y);
    	bool mustContinue= false;

		for(int i=*LocX; i<*LocX+templ.cols; i++){
			for(int j=*LocY; j<*LocY+templ.rows;j++){
				//printf("location : %d, %d\n",i ,j);
				if(comparePixels(img_display.at<Vec3b>(j,i), grayPix)){
					mustContinue = true;
					break;
				}
			}
			if (mustContinue){
				break;
			}
		}
		/*
		printf("Removing for : x:%d,y:%d = val:%f\n",*LocX,*LocY,	result.at<float>(*LocY,*LocX));
		for(int i=*LocX; i<*LocX+templ.cols; i++){
			for(int j=*LocY; j<*LocY+templ.rows;j++){
				//printf("location : %d, %d\n",i ,j);
				result.at<float>(j,i)=1.0;
			}
		}

		printf("Romoved : x:%d,y:%d = val:%f\n",*LocX,*LocY,	result.at<float>(*LocY,*LocX));
		if(mustContinue){
			//printf("Ended with:k = %d\n", k);
			break;
		}else{
			//printf("Done:k = %d\n", k);
		}

  		/// Show me what you got
  		rectangle( img_display, matchLoc, Point( matchLoc.x + templ.cols , matchLoc.y + templ.rows ), grayScalar, 2, 8, 0 );
  		//rectangle( result, matchLoc, Point( matchLoc.x + templ.cols , matchLoc.y + templ.rows ), Scalar::all(128), 2, 8, 0 );

  		matchTemplate( img_display, templ, result, match_method );
  	}

  	imshow( image_window, img_display );
  	imshow( result_window, result );

	/*
  	minMaxLoc( result, &minVal, &maxVal, &minLoc, &maxLoc, Mat() );

	matchLoc = minLoc;

  	rectangle( img_display, matchLoc, Point( matchLoc.x + templ.cols , matchLoc.y + templ.rows ), Scalar::all(128), 2, 8, 0 );
  	rectangle( result, matchLoc, Point( matchLoc.x + templ.cols , matchLoc.y + templ.rows ), Scalar::all(64), 2, 8, 0 );

}
  	*/