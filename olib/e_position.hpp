
#ifndef EPOSITION_H
#define EPOSITION_H

#include "string.h"
#include "BasicPoint.hpp"
#include "utils.hpp"

class e_position
{
public:
	e_position(BasicPoint topLeftP, BasicPoint bottomRightP);
	
	BasicPoint p1,p2;
	BasicPoint midPoint() const;
	std::string toString() const;	
	friend std::ostream& operator<< (std::ostream& outs, const e_position& obj);
};

	std::ostream& operator<< (std::ostream& outs, const e_position& obj);


#endif