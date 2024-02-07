package object

import "testing"

func TestStringHashKeys(t *testing.T) {
	hello1 := &String{Value: "Hello world"}
	hello2 := &String{Value: "Hello world"}
	diff1 := &String{Value: "Here is Johnny"}
	diff2 := &String{Value: "Here is Johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Fatalf("strings with the same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Fatalf("strings with different content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Fatalf("strings with different content have the same hash keys")
	}

}
