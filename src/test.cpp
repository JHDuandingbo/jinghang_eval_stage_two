#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"
#include <iostream>

using namespace rapidjson;

int main() {
	// 1. Parse a JSON string into DOM.
	const char* json = "{\"project\":\"rapidjson\",\"a\":{\"a\":122,\"b\":\"hello\"}, \"stars\":10}";
	Document d;
	d.SetObject();
	d.Parse(json);

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
	StringBuffer buffer;
	Writer<StringBuffer> writer(buffer);
	d.Accept(writer);

	std::cout<<d["a"]["a"].GetInt()<<std::endl;
	std::cout<<d["a"]["b"].GetString()<<std::endl;
	// Output {"project":"rapidjson","stars":11}
	std::cout << buffer.GetString() << std::endl;
	return 0;
}
