package sorting

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	testData := []string{"a", "z", "h", "e", "t", "u", "e", "y", "r", "h", "u", "i"}
	sorted := Split(testData)
	sortedAnswer := []string{"a", "e", "e", "h", "h", "i", "r", "t", "u", "u", "y", "z"}
	if len(sorted) != len(sortedAnswer) {
		t.Errorf("Expected %d got %d", len(sortedAnswer), len(sorted))
	}
	for i, v := range sorted {
		if v != sortedAnswer[i] {
			t.Errorf("Expected %s got %s", sortedAnswer[i], v)
		}
	}
}
