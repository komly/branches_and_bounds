package branches_and_bounds

import "testing"

func sliceEquals(a, b []int) bool{
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true

}

func mapKeysEquals(a, b map[int]int) bool{
	keysA := make([]int, 0, len(a))
	for k := range a {
		keysA = append(keysA, k)
	}

	keysB := make([]int, 0, len(b))
	for k := range b {
		keysB = append(keysB, k)
	}

	return sliceEquals(keysA, keysB)

}

func TestRemoveRowAndCell(t *testing.T) {
	data := [][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	m := MatrixFromData(data)
	m.removeRowAndCol(0, 0)

	mapKeysEquals(m.rows, map[int]int{
		1: 0,
		2: 1,
	})

	mapKeysEquals(m.cols, map[int]int{
		1: 0,
		2: 1,
	})
}

func TestBranchesAndBounds(t *testing.T) {
	BranchesAndBounds([][]float64{}, []int{})
}
