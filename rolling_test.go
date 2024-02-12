package num

import (
	"reflect"
	"testing"
)

func TestRolling(t *testing.T) {
	testSliceFloat := []float64{1, 2, 3, 4, 5, 6}
	//expected := []float64{1, 2, 3, 4, 5, 6}
	expected := [][]float64{
		{},
		{},
		{1, 2, 3},
		{2, 3, 4},
		{3, 4, 5},
		{4, 5, 6},
	}
	output := Rolling(testSliceFloat, 3)
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

	output = Rolling(testSliceFloat, 3)
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

	output = Rolling(testSliceFloat, []int{3, 3, 3, 3, 3, 3})
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

	output = Rolling(testSliceFloat, []float64{3, 3, 3, 3, 3, 3})
	if reflect.DeepEqual(expected, output) != true {
		t.Errorf("Got %v, want %v", output, expected)
	}

}
