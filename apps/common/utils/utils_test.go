package utils

import "testing"

func TestSliceToMap(t *testing.T) {

	type testStruct struct {
		Id   string
		Name string
	}

	s := []testStruct{
		{
			Id:   "1",
			Name: "1",
		},
		{
			Id:   "2",
			Name: "2",
		},
		{
			Id:   "3",
			Name: "3",
		},
	}

	m, err := SliceToMap[string, testStruct](s, "Id")

	if err != nil {
		t.Error(err)
	}

	if len(m) != 3 {
		t.Error("length not equal")
	}

	if m["1"].Name != "1" {
		t.Error("value not equal")
	}

	if m["2"].Name != "2" {
		t.Error("value not equal")
	}

	if m["3"].Name != "3" {
		t.Error("value not equal")
	}

}
