#ifndef POINTRELATION_H
#define POINTRELATION_H

#include "e_position.hpp"
#include <stdio.h>
#include "stdlib.h"
#include <vector>
#include "utils.hpp"

class PosRelation{
private:
	e_position freePosition;
	float centerAngle;
	int centerDistance;
public:
	PosRelation(e_position& p1,float p2, int p3);

	e_position getFreePosition();
	float getCenterAngle();
	int getCenterDistance();
	std::string toString() const;   
	friend std::ostream& operator<< (std::ostream& outs, const PosRelation& obj);
};

class PosRelations 
{
private:
		e_position parent;
        std::vector<PosRelation> east,south,west,north,southeast,northeast,southwest,northwest;
        bool addTo(vector<PosRelation> &v, PosRelation rel);
public:
        	//constructors
        PosRelations( e_position p1);
        //PosRelations( e_position p1,vector<PosRelation> p2);
        	// getters
        e_position getParent() const;
        std::vector<PosRelation> getEast() const;
		std::vector<PosRelation> getSouth() const;
		std::vector<PosRelation> getWest() const;
		std::vector<PosRelation> getNorth() const;
		std::vector<PosRelation> getSouthEast() const;
		std::vector<PosRelation> getNorthEast() const;
		std::vector<PosRelation> getSouthWest() const;
		std::vector<PosRelation> getNorthWest() const;
			//adders
        bool addToEast(PosRelation rel);
		bool addToSouth(PosRelation rel);
		bool addToWest(PosRelation rel);
		bool addToNorth(PosRelation rel);
        bool addToSouthEast(PosRelation rel);
		bool addToSouthWest(PosRelation rel);
		bool addToNorthWest(PosRelation rel);
		bool addToNorthEast(PosRelation rel);
			//methods
		int totalElements();
			//String Utils
		std::string toString() const;   
		friend std::ostream& operator<< (std::ostream& outs, const PosRelations& obj);
};

	std::ostream& operator<< (std::ostream& outs, const PosRelations& obj);
	std::ostream& operator<< (std::ostream& outs, const PosRelation& obj);

#endif