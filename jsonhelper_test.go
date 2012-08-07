package gou

import (
	"encoding/json"
	"testing"
)

//  go test -bench=".*" 
//  go test -run="(Util)"

func TestUtilJsonHelper(t *testing.T) {

	jh := make(JsonHelper)
	json.Unmarshal([]byte(`{
		"name":"string",
		"ints":[1,2,3,4],
		"int":1,
		"int64":1234567890,
		"MaxSize" : 1048576,
		"nested":{
			"nest":"string2",
			"int":2,
			"list":["value"],
			"nest2":{
				"test":"good"
			}
		},
		"nested2":[
			{"sub":2}
		]
	}`), &jh)

	Assert(jh.String("name") == "string", t, "should get string %s", jh.String("name"))
	Assert(jh.Int("int") == 1, t, "get int ")
	Assert(jh.Int("notint") == -1, t, "get non existent int = 0??? ")
	Assert(jh.Int("ints[0]") == 1, t, "get int from array %d", jh.Int("ints[0]"))
	Assert(jh.Int("ints[2]") == 3, t, "get int from array %d", jh.Int("ints[0]"))
	Assert(jh.Int64("int64") == 1234567890, t, "get int")
	Assert(jh.Int("nested.int") == 2, t, "get int")
	Assert(jh.String("nested.nest") == "string2", t, "should get string %s", jh.String("nested.nest"))
	Assert(jh.String("nested.nest2.test") == "good", t, "should get string %s", jh.String("nested.nest2.test"))
	Assert(jh.String("nested.list[0]") == "value", t, "get string from array")
	Assert(jh.Int("nested2[0].sub") == 2, t, "get int from obj in array %d", jh.Int("nested2[0].sub"))

	Assert(jh.Int("MaxSize") == 1048576, t, "get int, test capitalization? ")
}
