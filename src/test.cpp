#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include <iostream>

using namespace rapidjson;

int main() {
	// 1. Parse a JSON string into DOM.
	const char* bar = "{\"project\":\"rapidjson\",\"a\":{\"bar\":122,\"bbb\":\"hello\"}, \"stars\":10}";
	const char* foo = "{\"project\":\"rapidjson\",\"a\":{\"foo\":122,\"b\":\"hello\"}, \"stars\":10}";
	Document  fooD,barD;
	fooD.Parse(foo);
	barD.Parse(bar);



	//Value::AllocatorType allocator;
	//Value x(kObjectType);

//	x.AddMember("foo", 9999, x.GetAllocator());

	/*
	   Value foo;
	   foo.SetObject();
	   foo.AddMember("foofoo", 11111, d.GetAllocator());
	   d.AddMember("foo", foo, d.GetAllocator());

	// 2. Modify it by DOM.
	Value& s = d["stars"];
	s.SetInt(s.GetInt() + 1);

	//d.AddMember("test", "ok");
	// 3. Stringify the DOM
	const char *ccc="cccccccccccc";
	d.AddMember("abc", Value("").SetString(ccc, 100), d.GetAllocator());
	//	d.AddMember("ccc", ccc, d.GetAllocator());
	 */
	barD.RemoveMember("a");
	barD.AddMember("a", fooD["a"],barD.GetAllocator());
	StringBuffer buffer;
	Writer<StringBuffer> writer(buffer);
	barD.Accept(writer);
//	std::cout << buffer.GetString() << std::endl;
	//x.RemoveMember("foo");

	//d.AddMember("a", dd["a"], d.GetAllocator());
	//	dd["request"].Swap(d["request"]);
	//d.SetObject("ccc", dd["a"], d.GetAllocator());
	std::cout << buffer.GetString() << std::endl;
	return 0;
}
