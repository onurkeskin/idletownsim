
#ifndef BASICPONT_H
#define BASICPONT_H

#include "utils.hpp"
#include <iostream>
#include <math.h>

using namespace std;

// Class to represent BasicPoints.
class BasicPoint {
private:
        double xval, yval;
public:
        // Constructor uses default arguments to allow calling with zero, one,
        // or two values.
        BasicPoint(double x = 0.0, double y = 0.0);
        // Extractors.
        double x() const;
        double y() const;

        // Distance to another BasicPoint.  Pythagorean thm.
        double dist(BasicPoint other);
        double fastDist(BasicPoint other);
        double angle(BasicPoint other);
        // Add or subtract two BasicPoints.
        BasicPoint add(BasicPoint b);
        BasicPoint sub(BasicPoint b);
        // Move the existing BasicPoint.
        void move(double a, double b);

        // Print the BasicPoint on the stream.  The class ostream is a base class
        // for output streams of various types.
        void print(ostream &strm);

        std::string toString() const;   
        friend std::ostream& operator<< (std::ostream& outs, const BasicPoint& obj);
};

        std::ostream& operator<< (std::ostream& outs, const BasicPoint& obj);

#endif