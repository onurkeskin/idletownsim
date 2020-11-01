#include "e_position.hpp"

using namespace std;

e_position::e_position(BasicPoint topLeftP, BasicPoint bottomRightP){
	p1 = topLeftP;
	p2 = bottomRightP;
}

BasicPoint e_position::midPoint() const{
	return BasicPoint((p1.x() + p2.x())/2, (p1.y()+p2.y())/2);
}

std::string e_position::toString() const{
	string str = "Pos: " + IntToString(p1.x()) + "," + IntToString(p1.y()) + " to " + IntToString(p2.x()) + "," + IntToString(p2.y());
	return str;
}
	
std::ostream& operator<< (std::ostream& outs, const e_position& obj) {
	return outs << obj.toString();
}
