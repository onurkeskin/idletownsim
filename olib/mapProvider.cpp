#include "opencv2/opencv.hpp"
#include "opencv/cv.h"
#include <string>
#include <stdio.h>
#include "stdlib.h"
#include "mapProvider.hpp"
#include "cWrapper.h"
#include <vector>
#include "EmptyPosRelation.hpp"
#include <iterator>

using namespace cv;
using namespace std;

bool comparePixels(Vec3b p1, Vec3b p2){
	if (p1[0] == p2[0] && p1[1] == p2[1] && p1[2] == p2[2]){
		return true;
	}
	return false;
}

void printMat(Mat mat){
	int c = mat.channels();

	for(int i=0; i<mat.cols; i++){
		for(int j=0; j<mat.rows;j++){
			printf("%d ",mat.at<uchar>(j,i));
		}
		printf("\n");
	}
}

bool hardChangeColorsToBlackAndWhite(Mat &img){
	int c = img.channels();
	if(c != 1) return false;

	for(int i=0; i<img.cols; i++){
		for(int j=0; j<img.rows;j++){
			if(img.at<uchar>(j,i) <=128) {
				img.at<uchar>(j,i) = 0;
			} else {
				img.at<uchar>(j,i) = 255;
			}
		}
	}
	return true;
}

int parseMapW(Mat image,Mat &result,Mat &displayResult, vector<e_position> &positions, Mat templ = Mat(40, 40, CV_8U, Scalar(0)), int matchFunction = CV_TM_SQDIFF,int maxMatches = 1000,uchar highlighter = (uchar)50){
	Mat toUse;
	if (image.channels()!= 1) {
		cvtColor(image,toUse,CV_BGR2GRAY,1);
	}
	hardChangeColorsToBlackAndWhite(toUse);

	toUse.copyTo( displayResult );
  	//displayResult = image.clone();
	int result_cols =  toUse.cols - templ.cols + 1;
	int result_rows = toUse.rows - templ.rows + 1;

	(result).create( result_rows, result_cols, CV_8U );
	matchTemplate( toUse, templ, result, matchFunction );

	double minVal; double maxVal; Point minLoc; Point maxLoc;
	Point matchLoc;

	int k;
	for(k=1;k <= maxMatches;k++) {
		minMaxLoc( result, &minVal, &maxVal, &minLoc, &maxLoc, Mat() );

		if( matchFunction  == CV_TM_SQDIFF || matchFunction == CV_TM_SQDIFF_NORMED )
			{ matchLoc = minLoc; }
		else
			{ matchLoc = maxLoc; }

		int LocX = (matchLoc.x);
		int LocY = (matchLoc.y);

		bool mustContinue= false;

		for(int i=LocX; i<LocX+templ.cols; i++){
			for(int j=LocY; j<LocY+templ.rows;j++){
				//printf("location : %d, %d\n",i ,j);
				if((displayResult).at<uchar>(j,i) == highlighter){
					//printf("wtf : %d, %d\n",(displayResult).at<uchar>(j,i) ,highlighter);
					mustContinue = true;
					break;
				}
			}
			if (mustContinue){
				break;
			}
		}

		if(mustContinue){
			//printf("Ended with:k = %d\n", k-1);
			break;
		}else{
			//printf("Done:k = %d\n", k);
		}

		rectangle( displayResult, matchLoc, Point( matchLoc.x + templ.cols , matchLoc.y + templ.rows ), highlighter, 1, 8, 0 );
		putText(displayResult, IntToString(k), Point( (matchLoc.x + templ.cols/2-10), (matchLoc.y+ templ.rows/2)) , FONT_HERSHEY_PLAIN, 0.8,highlighter, 1, 8);

		e_position pos(BasicPoint(LocX,LocY),BasicPoint(LocX+templ.cols,matchLoc.y + templ.rows));
		positions.push_back(pos);
  		//rectangle( result, matchLoc, Point( matchLoc.x + templ.cols , matchLoc.y + templ.rows ), Scalar::all(128), 2, 8, 0 );

		matchTemplate( displayResult, templ, result, matchFunction );
	}

	return 0;
}

vector<BasicPoint> ExtractMidPointsFromRect(vector<e_position> rects){
	vector <BasicPoint> toRet = vector<BasicPoint>();
	for(std::vector<e_position>::iterator it = rects.begin(); it != rects.end(); ++it) {
		toRet.push_back(it->midPoint());
	}

	return toRet;
}

int countRelations(vector<PosRelations> relations){
	int TotalCount = 0;
	for(vector<PosRelations>::iterator it = relations.begin(); it != relations.end(); it++){
		PosRelations rel = *it;
		TotalCount+=rel.totalElements();
	}
	return TotalCount;
}

vector<PosRelations> AnalyzeRects(vector<e_position> rects,int maxDist, float piCoefficient){
	vector<BasicPoint> points = ExtractMidPointsFromRect(rects);
	int c = 1;
	vector<PosRelations> allRelations = vector<PosRelations>();
	for(std::vector<BasicPoint>::iterator it = points.begin(); it != points.end(); ++it) {
		BasicPoint curPoint = *it;

		PosRelations relations = PosRelations(rects[c-1]);
		c++;
		//vector<e_position*> closestnorth = vector<e_position*>(), closesteast = vector<e_position*>(), closestsouth = vector<e_position*>(), closestwest = vector<e_position*>(),diagonals = vector<e_position*>();
		int eastDist = INT_MAX, southDist = INT_MAX, westDist = INT_MAX, northDist = INT_MAX,northEastDist = INT_MAX,northWestDist = INT_MAX,southEastDist = INT_MAX,southWestDist = INT_MAX;
		for(std::vector<int>::size_type i = 0; i != points.size(); i++) {
			if(&(*it) == &(points[i])) { continue;}

			int distance = points[i].fastDist(*it);
			double angle = points[i].angle(*it);
			//cout << angle << " " << piCoefficient << endl;

			PosRelation rel = PosRelation((rects[i]),angle,distance);
			if (angle > -piCoefficient && angle < +piCoefficient){
				if(distance > maxDist) continue;
				//if(distance > eastDist) continue;
				eastDist = distance;
				relations.addToEast(rel);
			} else if (angle > M_PI/2-piCoefficient && angle < M_PI/2+piCoefficient ) {
				if(distance > maxDist) continue;
				//if(distance > southDist) continue;
				southDist = distance;
				relations.addToSouth(rel);
			} else if (angle > M_PI-piCoefficient && angle < M_PI+piCoefficient) {
				if(distance > maxDist) continue;
				//if(distance > westDist) continue;
				westDist = distance;
				relations.addToWest(rel);
			} else if (angle < -(M_PI/2-piCoefficient )&& angle > -(M_PI/2+piCoefficient) ) {
				if(distance > maxDist) continue;
				//if(distance > northDist) continue;
				northDist = distance;
				relations.addToNorth(rel);
			} else if (angle < -(M_PI-piCoefficient) && angle > -(M_PI+piCoefficient) ) {
				if(distance > maxDist) continue;
				//if(distance > westDist) continue;
				westDist = distance;
				relations.addToWest(rel);
			} else if(angle > (M_PI/4-piCoefficient) && angle < (M_PI/4+piCoefficient) ){
			if(distance > maxDist) continue;
				southEastDist = distance;
				relations.addToSouthEast(rel);
			} else if(angle > (3*M_PI/4-piCoefficient) && angle < (3*M_PI/4+piCoefficient) ){
				if(distance > maxDist) continue;
				southWestDist = distance;
				relations.addToSouthWest(rel);
			} else if(angle < -(M_PI/4-piCoefficient) && angle > -(M_PI/4+piCoefficient) ){
				if(distance > maxDist) continue;
				northEastDist = distance;
				relations.addToNorthEast(rel);
			} else if(angle < -(3*M_PI/4-piCoefficient) && angle > -(3*M_PI/4+piCoefficient) ){
				if(distance > maxDist) continue;
				northWestDist = distance;
				relations.addToNorthWest(rel);
			}
			//cout << c-1 << ": " << angle << " " << distance << endl;
		}
		//cout << c-1 << ": " << relations << endl;
		allRelations.push_back(relations);
	}
	/*
	int k = 1;
	for(std::vector<PosRelations>::iterator printit = allRelations.begin(); printit != allRelations.end(); ++printit) {
		cout << k++ << ":" << (*printit) << endl;
	}
	

	int relationTotal = countRelations(allRelations);
	cout<<relationTotal<<endl;
	*/

	return allRelations;
}


cEmptyPosRelation formCRelation(PosRelation rel){
	e_position rect = rel.getFreePosition();
	int x1 = rect.p1.x(), y1 = rect.p1.y() , x2 = rect.p2.x(), y2 = rect.p2.y();
	cBasicPoint point1{x1,y1};
	cBasicPoint point2{x2,y2};

	cEPosition position{point1, point2};
	cEmptyPosRelation toRet{position,rel.getCenterAngle(),rel.getCenterDistance()};
	//cout << x1 << " " << x2 << " " << y1 << " " << y2 << " " <<endl;
	//cout << rel << endl;
	return toRet;
}

void printCEPosition(cEPosition rect){
    cout << "x1:"<< rect.p1.xval <<" y1:" << rect.p1.yval << " x2:"<< rect.p2.xval <<" y2:" << rect.p2.yval << endl;
}


int testCExport(cExports exp){

    cout << "count:" << exp.count <<endl;
    for(int i = 0; i<exp.count; i++){
        cEmptyPosRelations cur = (exp.exports[i]);
        cEPosition parent = cur.parent;
        
        cEmptyPosRelation *east = cur.east;
        cEmptyPosRelation *west = cur.west;
        cEmptyPosRelation *south = cur.south;
        cEmptyPosRelation *north = cur.north;
        cEmptyPosRelation *southeast = cur.southeast;
        cEmptyPosRelation *southwest = cur.southwest;
        cEmptyPosRelation *northeast = cur.northeast;
        cEmptyPosRelation *northwest = cur.northwest;

        cout << "----------------------------"<< i+1 << "----------------------------" << endl;
        cout << "For Point: ";
        printCEPosition(parent);

        cout << "easts: " << endl;
        for(int r = 0; r<cur.eastSize; r++){
            cEmptyPosRelation va = east[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "souths: " << endl;
        for(int r = 0; r<cur.southSize; r++){
            cEmptyPosRelation va = south[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "wests: " << endl;
        for(int r = 0; r<cur.westSize; r++){
            cEmptyPosRelation va = west[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "norths: " << endl;
        for(int r = 0; r<cur.northSize; r++){
            cEmptyPosRelation va = north[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "southeasts: " << endl;
        for(int r = 0; r<cur.southeastSize; r++){
            cEmptyPosRelation va = southeast[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "southwests: " << endl;
        for(int r = 0; r<cur.southwestSize; r++){
            cEmptyPosRelation va = southwest[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "northwests: " << endl;
        for(int r = 0; r<cur.northwestSize; r++){
            cEmptyPosRelation va = northwest[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
        cout << "northeasts: " << endl;
        for(int r = 0; r<cur.northeastSize; r++){
            cEmptyPosRelation va = northeast[r];
            cout << "Pos: ";
            printCEPosition(va.freePosition);
            cout << " Angle: " << va.centerAngle << " Dist: " << va.centerDistance << " " << endl;
        }
            cout <<endl;
    }


}

cEmptyPosRelations* formCPackage(vector<PosRelations> &rels){
	cEmptyPosRelations* cRels = (cEmptyPosRelations*)malloc(sizeof(cEmptyPosRelations)*rels.size());
	int cRelCount = 0;
	for(std::vector<PosRelations>::iterator relIt = rels.begin(); relIt != rels.end(); ++relIt) {

		e_position rect = (*relIt).getParent();
		int x1 = rect.p1.x(), y1 = rect.p1.y() , x2 = rect.p2.x(), y2 = rect.p2.y();
		cBasicPoint* parentPoint1 = (cBasicPoint*)malloc(sizeof(cBasicPoint));
		parentPoint1->xval = x1; parentPoint1->yval = y1;
		cBasicPoint* parentPoint2 = (cBasicPoint*)malloc(sizeof(cBasicPoint));
		parentPoint2->xval = x2; parentPoint2->yval = y2;
		cEPosition* parentPos = (cEPosition*)malloc(sizeof(cEPosition));
		parentPos->p1 = *parentPoint1; parentPos->p2 = *parentPoint2;

		vector<PosRelation> east = (*relIt).getEast();
		vector<PosRelation> south = (*relIt).getSouth();
		vector<PosRelation> west = (*relIt).getWest();
		vector<PosRelation> north = (*relIt).getNorth();
		vector<PosRelation> southeast = (*relIt).getSouthEast();
		vector<PosRelation>	northeast = (*relIt).getNorthEast();
		vector<PosRelation> southwest = (*relIt).getSouthWest();
		vector<PosRelation> northwest = (*relIt).getNorthWest();

		cRels[cRelCount].parent = *parentPos;
		cRels[cRelCount].east = /*cEmptyPosRelation[east.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*east.size());
		cRels[cRelCount].south = /*cEmptyPosRelation[south.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*south.size());
		cRels[cRelCount].west = /*cEmptyPosRelation[west.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*west.size());
		cRels[cRelCount].north = /*cEmptyPosRelation[north.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*north.size());
		cRels[cRelCount].southeast = /*cEmptyPosRelation[others.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*southeast.size());
		cRels[cRelCount].northeast = /*cEmptyPosRelation[others.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*northeast.size());
		cRels[cRelCount].southwest = /*cEmptyPosRelation[others.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*southwest.size());
		cRels[cRelCount].northwest = /*cEmptyPosRelation[others.size()];*/ (cEmptyPosRelation*)malloc(sizeof(cEmptyPosRelation)*northwest.size());

		cRels[cRelCount].eastSize = east.size();;
		cRels[cRelCount].southSize = south.size();;
		cRels[cRelCount].westSize = west.size();;
		cRels[cRelCount].northSize = north.size();;
		cRels[cRelCount].southeastSize = southeast.size();
		cRels[cRelCount].northeastSize = northeast.size();
		cRels[cRelCount].southwestSize = southwest.size();
		cRels[cRelCount].northwestSize = northwest.size();

		int cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = east.begin(); it != east.end(); ++it) {
			//cout << *it << endl;
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].east[cInnerRel] = cRel;
			//cout<<cRel.centerDistance<< " a " << cRel.centerAngle << " p: x:" << cRel.freePosition.p1.xval << " y:" << cRel.freePosition.p1.yval << " x:" << cRel.freePosition.p2.xval <<" y:" <<cRel.freePosition.p2.yval <<endl; 
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = south.begin(); it != south.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].south[cInnerRel] = cRel;
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = west.begin(); it != west.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].west[cInnerRel] = cRel;
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = north.begin(); it != north.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].north[cInnerRel] = cRel;
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = southeast.begin(); it != southeast.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].southeast[cInnerRel] = cRel;
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = northeast.begin(); it != northeast.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].northeast[cInnerRel] = cRel;
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = southwest.begin(); it != southwest.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].southwest[cInnerRel] = cRel;
			cInnerRel++;
		}
		cInnerRel = 0;
		for(std::vector<PosRelation>::iterator it = northwest.begin(); it != northwest.end(); ++it) {
			cEmptyPosRelation cRel;
			cRel = formCRelation(*it);
			cRels[cRelCount].northwest[cInnerRel] = cRel;
			cInnerRel++;
		}
		cRelCount++;
	}

	return cRels;
}


cExports basic(unsigned char* img,int length) {
	//Mat rawData  =  Mat( 1, length, CV_8UC1, img );
	//std::vector<unsigned char> data;
	//std::copy(std::istream_iterator<unsigned char>(img), std::istream_iterator<unsigned char>(), std::back_inserter(data));
	std::vector<unsigned char> data(img, img + length);
	//Mat matrixJprg = imdecode(Mat(data), 1);
	Mat image = imdecode(Mat(data), 1);//imdecode((InputArray)img,CV_LOAD_IMAGE_UNCHANGED);

/*
	for (std::vector<unsigned char>::const_iterator i = data.begin(); i != data.end(); ++i){
    	std::cout << +(*i) << ' ';
	}
	cout<<endl;
*/

	Mat croppedImg = image(Rect(0,0,400,380));

	cExports toRet;
	char* image_window = "Reflected Image";
	char* result_window = "Result window";
	//char* final_window = "Original";

	//namedWindow( image_window, CV_WINDOW_AUTOSIZE );
	//namedWindow( result_window, CV_WINDOW_AUTOSIZE );
	
	/*
    char* trackbar_label = "Method: \n 0: SQDIFF \n 1: SQDIFF NORMED \n 2: TM CCORR \n 3: TM CCORR NORMED \n 4: TM COEFF \n 5: TM COEFF NORMED";
    char* count_label = "Count";
    createTrackbar( trackbar_label, image_window, &match_method, max_Trackbar, parseMap );
    createTrackbar( count_label, image_window, &spaceCount, maxSpaceCount, parseMap );
    parseMap(0,0);

	*/

	Mat result, displayResult;
	vector<e_position> vec = vector<e_position>();
	parseMapW(croppedImg, result, displayResult,vec,Mat(30, 30, CV_8U, Scalar(0)));
    //cvtColor(croppedImg,result,CV_BGR2GRAY,1);
    //printf("channel count : %d \n",result.channels());
    //printf("result : %d \n", hardChangeColorsToBlackAndWhite(result));
    //printMat(result);
	//imshow( result_window, result );
	//imshow( image_window, displayResult );

 	/*
 	for (vector<e_position>::const_iterator i = vec.begin(); i != vec.end(); ++i){
    	cout << *i << endl;
 	}
 	*/

	vector<PosRelations> rels = AnalyzeRects(vec,2500,M_PI/6);

	/*
	int asda = 0;
 	for (vector<PosRelations>::const_iterator i = (rels).begin(); i != (rels).end(); ++i){
 		cout << "cur: " << asda << endl;
 		PosRelations north = (*i);
    	cout << north << endl;
    	asda++;
 	}
 	*/
	std::vector<unsigned char> vectImg;
	imencode(".png",croppedImg,vectImg);



	cEmptyPosRelations *cRels = formCPackage(rels);
	toRet.count = rels.size();
	toRet.exports = cRels;

	toRet.imgsize = vectImg.size();

	unsigned char* retImg = (unsigned char*) malloc(sizeof(unsigned char) * vectImg.size());
	std::copy(vectImg.begin(), vectImg.end(), retImg);

	toRet.modimg = &retImg[0];
	//testCExport(toRet);
	
	/*
	cout<<endl;
	cout<<"End Image"<<endl;
	for (std::vector<unsigned char>::const_iterator i = vectImg.begin(); i != vectImg.end(); ++i){
    	std::cout << +(*i) << ' ';
	}
	cout<<endl;
*/
	return toRet;
}
