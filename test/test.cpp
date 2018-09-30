// condition_variable example
#include <iostream>           // std::cout
#include <thread>             // std::thread
#include <mutex>              // std::mutex, std::unique_lock
#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"



using namespace rapidjson;
using namespace std;




const char *  test(Value  &v){




	StringBuffer buffer;
	Writer<StringBuffer> writer(buffer);
	v.Accept(writer);
	const char * str = buffer.GetString();
	cout <<str<<endl;

	return strdup(str);
	//return "";
}
	
int main ()
{

	Document d1;
	Document::AllocatorType & a = d1.GetAllocator();
	d1.SetObject();
	d1.AddMember("aa",111, d1.GetAllocator());
	

	Value b;
	b.SetArray();
	b.PushBack(1, a);
	b.PushBack(1, a);
	b.PushBack("heell", a);
/*
	Document::AllocatorType& a = d1.GetAllocator();

	d1.SetArray().PushBack(1, a).PushBack(2, a);

	StringBuffer buffer;
	Writer<StringBuffer> writer(buffer);
	d1.Accept(writer);

	cout << buffer.GetString()<<endl;
	const char * aa = "1";;
	double bb = 2.2;
*/
	cout<< test(b) <<endl;
	cout<< test(d1) <<endl;
}
