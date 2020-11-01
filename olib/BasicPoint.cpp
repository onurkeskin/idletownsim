#include "BasicPoint.hpp"
// Constructor uses default arguments to allow calling with zero, one,
// or two values.


BasicPoint::BasicPoint(double x /*= 0.0*/, double y /*= 0.0*/) {
xval = x;
yval = y;
}

                // Extractors.
double BasicPoint::x() const { return xval; }
double BasicPoint::y() const{ return yval; }

                // Distance to another BasicPoint.  Pythagorean thm.
double BasicPoint::dist(BasicPoint other) {
        double xd = xval - other.xval;
        double yd = yval - other.yval;
        return sqrt(xd*xd + yd*yd);
}

double BasicPoint::fastDist(BasicPoint other){
        double xd = xval - other.xval;
        double yd = yval - other.yval;
        return (xd*xd + yd*yd);
}

double BasicPoint::angle(BasicPoint other){
        return atan2(yval - other.yval, xval - other.xval);
}

                // Add or subtract two BasicPoints.
BasicPoint BasicPoint::add(BasicPoint b)
{
        return BasicPoint(xval + b.xval, yval + b.yval);
}
BasicPoint BasicPoint::sub(BasicPoint b)
{
        return BasicPoint(xval - b.xval, yval - b.yval);
}

                // Move the existing BasicPoint.
void BasicPoint::move(double a, double b)
{
        xval += a;
        yval += b;
}

                // Print the BasicPoint on the stream.  The class ostream is a base class
                // for output streams of various types.
void BasicPoint::print(ostream &strm)
{
        strm << "(" << xval << "," << yval << ")";
}

std::string BasicPoint::toString() const{
        string str = IntToString(x()) + "," + IntToString(y());
        return str;
}

std::ostream& operator<< (std::ostream& outs, const BasicPoint& obj) {
        return outs << obj.toString();
}
