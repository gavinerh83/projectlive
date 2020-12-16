//Package sorting implements the sorting algorithm for merge sort
package sorting

//Split splits the slice of string into 2 and recursively do so until the length slice becomes 1
func Split(num []string) []string {
	if len(num) == 1 {
		return num
	}
	middle := len(num) / 2
	left := num[:middle]
	right := num[middle:]

	return merge(Split(left), Split(right))
}

//merge combines the previously separated slice in order
func merge(left, right []string) []string {
	var result []string
	leftindex := 0
	rightindex := 0
	for leftindex < len(left) && rightindex < len(right) { //stop the merge when either of the 2 slices finishes
		if left[leftindex] > right[rightindex] {
			result = append(result, right[rightindex])
			rightindex++
		} else {
			result = append(result, left[leftindex])
			leftindex++
		}
	}
	//if the left or right index still has elements
	for leftindex < len(left) {
		result = append(result, left[leftindex])
		leftindex++
	}
	for rightindex < len(right) {
		result = append(result, right[rightindex])
		rightindex++
	}
	return result
}
