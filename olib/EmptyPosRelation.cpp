#include "EmptyPosRelation.hpp"

using namespace std;

PosRelation::PosRelation(e_position &p1,float p2, int p3)
:freePosition(p1)
,centerAngle(p2)
,centerDistance(p3)
{

}

e_position PosRelation::getFreePosition(){
	return freePosition;
}

float PosRelation::getCenterAngle(){
 	return centerAngle;
}
int PosRelation::getCenterDistance(){
	return centerDistance;
}

PosRelations::PosRelations( e_position p1)
:parent(p1)
{
	east = vector<PosRelation>();
	south = vector<PosRelation>();
	west = vector<PosRelation>();
	north = vector<PosRelation>();
	southeast = vector<PosRelation>();
	northeast = vector<PosRelation>();
	southwest = vector<PosRelation>();
	northwest = vector<PosRelation>();
}

vector<PosRelation> PosRelations::getEast() const{
	return east;
}
vector<PosRelation> PosRelations::getSouth() const{
	return south;
}
vector<PosRelation> PosRelations::getWest() const{
	return west;
}
vector<PosRelation> PosRelations::getNorth() const{
	return north;
}
vector<PosRelation> PosRelations::getSouthEast() const{
	return southeast;
}	
vector<PosRelation> PosRelations::getNorthEast() const{
	return northeast;
}	
vector<PosRelation> PosRelations::getSouthWest() const{
	return southwest;
}	
vector<PosRelation> PosRelations::getNorthWest() const{
	return northwest;
}	
e_position PosRelations::getParent() const {
	return parent;
}


bool PosRelations::addToEast(PosRelation rel){
	addTo(east,rel);
}
bool PosRelations::addToSouth(PosRelation rel){
	addTo(south,rel);
}
bool PosRelations::addToWest(PosRelation rel){
	addTo(west,rel);
}
bool PosRelations::addToNorth(PosRelation rel){
	addTo(north,rel);
}
bool PosRelations::addToSouthEast(PosRelation rel){
	addTo(southeast,rel);
}
bool PosRelations::addToNorthEast(PosRelation rel){
	addTo(northeast,rel);
}
bool PosRelations::addToSouthWest(PosRelation rel){
	addTo(southwest,rel);
}
bool PosRelations::addToNorthWest(PosRelation rel){
	addTo(northwest,rel);
}

bool PosRelations::addTo(vector<PosRelation> &v, PosRelation rel){
	v.push_back(rel);
}

int PosRelations::totalElements(){
	return east.size()+south.size()+west.size()+north.size()+southeast.size()+southwest.size()+northwest.size()+northeast.size();
}


std::string PosRelations::toString() const{
	string toRet = "";
	toRet+= "Parent: " + parent.toString();

	if(east.size() > 0){
		toRet+="\neasts: ";
		for(std::vector<int>::size_type i = 0; i != east.size(); i++) {
			toRet+=(east)[i].toString();
			toRet+= "|";
		}
	}
	if(south.size() > 0){
		toRet+="\nsouth: ";
		for(std::vector<int>::size_type i = 0; i != south.size(); i++) {
			toRet+=(south)[i].toString();
			toRet+= "|";
		}
	}
	if(west.size() > 0){
		toRet+="\nwest: ";
		for(std::vector<int>::size_type i = 0; i != west.size(); i++) {
			toRet+=(west)[i].toString();
			toRet+= "|";
		}
	}
	if(north.size() > 0){
		toRet+="\nnorth: ";
		for(std::vector<int>::size_type i = 0; i != north.size(); i++) {
			toRet+=(north)[i].toString();
			toRet+= "|";
		}
	}
	if(southeast.size() > 0){
		toRet+="\neasts: ";
		for(std::vector<int>::size_type i = 0; i != southeast.size(); i++) {
			toRet+=(southeast)[i].toString();
			toRet+= "|";
		}
	}
	if(southwest.size() > 0){
		toRet+="\nsouth: ";
		for(std::vector<int>::size_type i = 0; i != southwest.size(); i++) {
			toRet+=(southwest)[i].toString();
			toRet+= "|";
		}
	}
	if(northwest.size() > 0){
		toRet+="\nwest: ";
		for(std::vector<int>::size_type i = 0; i != northwest.size(); i++) {
			toRet+=(northwest)[i].toString();
			toRet+= "|";
		}
	}
	if(northeast.size() > 0){
		toRet+="\nnorth: ";
		for(std::vector<int>::size_type i = 0; i != northeast.size(); i++) {
			toRet+=(northeast)[i].toString();
			toRet+= "|";
		}
	}
	//toRet+="\n";

	return toRet;
}
std::string PosRelation::toString() const{
	return "Free Pos: " +freePosition.toString() + " Angle: " + FloatToString(centerAngle) + " Distance: " + IntToString(centerDistance);
}
std::ostream& operator<< (std::ostream& outs, const PosRelations& obj){
	return outs << obj.toString();
}
std::ostream& operator<< (std::ostream& outs, const PosRelation& obj){
	return outs << obj.toString();
}