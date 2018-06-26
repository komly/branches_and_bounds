package branches_and_bounds

import (
"fmt"
"math"
"sort"
"strings"
	"container/heap"
	"log"
)

type Matrix struct {
	cols map[int]int
	rows map[int]int
	m    [][]float64
}

type pathElem struct {
	r int
	c int
}

func MatrixFromData(data [][]float64) *Matrix {
	newMatrix := make([][]float64, 0)
	for _, row := range data {
		newRow := make([]float64, len(row))
		copy(newRow, row)
		newMatrix = append(newMatrix, newRow)
	}

	for i := 0; i < len(newMatrix); i++ {
		for j := 0; j < len(newMatrix[i]); j++ {
			if i == j {
				newMatrix[i][j] = math.MaxFloat64
			}
		}
	}

	rows := make(map[int]int)
	cols := make(map[int]int)

	for i := 0; i < len(data); i++ {
		rows[i] = i
	}

	for i := 0; i < len(data[0]); i++ {
		cols[i] = i
	}
	return &Matrix{
		cols: cols,
		rows: rows,
		m:    newMatrix,
	}
}

func (m *Matrix) get(r, c int) float64 {
	i := m.rows[r]
	j := m.cols[c]

	return m.m[i][j]
}

func (m *Matrix) set(r, c int, val float64) {
	i := m.rows[r]
	j := m.cols[c]
	m.m[i][j] = val
}

func (m *Matrix) debugPrint() string {
	res := make([]string, 0)

	rows := make([]int, 0)
	for r := range m.rows {
		rows = append(rows, r)
	}
	sort.Ints(rows)

	cols := make([]int, 0)
	for c := range m.cols {
		cols = append(cols, c)
	}
	sort.Ints(cols)

	for _, r := range rows {
		srow := make([]string, 0)
		for _, c := range cols {
			val := m.get(r, c)
			if val == math.MaxFloat64 {
				srow = append(srow, "M")
			} else {
				srow = append(srow, fmt.Sprintf("%.2f", val))
			}
		}
		res = append(res, strings.Join(srow, "\t"))

	}

	return strings.Join(res, "\n")
}

func (m *Matrix) getMinForRows() map[int]float64 {
	minForRows := make(map[int]float64)

	for r := range m.rows {
		minForRows[r] = math.MaxFloat64
	}
	for r := range m.rows {
		minForRows[r] = math.MaxFloat64
		for c := range m.cols {
			val := m.get(r, c)
			if val < minForRows[r] {
				minForRows[r] = val
			}
		}
	}
	return minForRows
}

func (m *Matrix) getPenalty(zr, zc int) float64 {
	minRow := math.MaxFloat64
	minCol := math.MaxFloat64

	for r := range m.rows {
		if r == zr {
			continue
		}
		val := m.get(r, zc)
		if val < minRow {
			minRow = val
		}
	}
	for c := range m.cols {
		if c == zc {
			continue
		}

		val := m.get(zr, c)
		if val < minCol {
			minCol = val
		}
	}
	return minRow + minCol
}

func (m *Matrix) getMinForCols() map[int]float64 {
	minForCols := make(map[int]float64)

	for c := range m.cols {
		minForCols[c] = math.MaxFloat64

	}
	for r := range m.rows {
		for c := range m.cols {
			val := m.get(r, c)
			if val < minForCols[c] {
				minForCols[c] = val
			}
		}
	}
	return minForCols
}

func (m *Matrix) clone() *Matrix {
	data := make([][]float64, 0)
	for _, row := range m.m {
		newRow := make([]float64, 0)
		for _, c := range row {
			newRow = append(newRow, c)
		}
		data = append(data, newRow)
	}

	newRows := make(map[int]int)
	for r, v := range m.rows {
		newRows[r] = v
	}

	newCols := make(map[int]int)
	for c, v := range m.cols {
		newCols[c] = v
	}


	marix := &Matrix{
		m:    data,
		rows: newRows,
		cols: newCols,
	}

	return marix
}

func (m *Matrix) removeRowAndCol(zr, zc int) *Matrix {
	if _, ok := m.rows[zr]; !ok {
		panic("can't delete not existing row")
	}

	if _, ok := m.cols[zc]; !ok {
		panic("can't delete not existing row")
	}

	rows := make([]int, 0)
	for r := range m.rows {
		rows = append(rows, r)
	}
	sort.Ints(rows)

	cols := make([]int, 0)
	for c := range m.cols {
		cols = append(cols, c)
	}
	sort.Ints(cols)

	data := make([][]float64, 0)
	for _, r := range rows {
		if r == zr {
			continue
		}
		newRow := make([]float64, 0)
		for _, c := range cols {
			if c == zc {
				continue
			}
			newRow = append(newRow, m.get(r, c))
		}
		data = append(data, newRow)
	}

	newRows := make(map[int]int)
	for r, v := range m.rows {
		if r == zr {
			continue
		}
		if r > zr {
			newRows[r] = v - 1


		} else {
			newRows[r] = v
		}
	}

	newCols := make(map[int]int)
	for c, v := range m.cols {
		if c == zc {
			continue
		}
		if c > zc {
			newCols[c] = v - 1
		} else {
			newCols[c] = v
		}
	}


	matrix := &Matrix{
		m:    data,
		rows: newRows,
		cols: newCols,
	}

	return matrix
}

func (m *Matrix) fmap(f func(float64, int, int) float64) {
	for r := range m.rows {
		for c := range m.cols {
			v := m.get(r, c)
			m.set(r, c, f(v, r, c))
		}
	}
}

func (m *Matrix) each(f func(float64, int, int)) {
	for r := range m.rows {
		for c := range m.cols {
			v := m.get(r, c)
			f(v, r, c)
		}
	}
}

type Solution struct {
	parent     *Solution
	matrix     *Matrix
	minBound   float64
	zr, zc     int
	maxPenalty float64
	index    int
	path []pathElem
}

func sum(m map[int]float64) float64 {
	res := 0.0
	for _, val := range m {
		if val == math.MaxFloat64 {
			continue
		}
		res += val
	}
	return res
}

func (s *Solution) reduce() bool {
	minForRows := s.matrix.getMinForRows()
	for _, min := range minForRows {
		if min == math.MaxFloat64 {
			return false
		}
	}

	s.matrix.fmap(func(v float64, r, c int) float64 {
		if v != math.MaxFloat64 {
			return v - minForRows[r]
		}
		return v
	})

	minForCols := s.matrix.getMinForCols()
	for _, min := range minForRows {
		if min == math.MaxFloat64 {
			return false
		}
	}
	for _, min := range minForCols {
		if min == math.MaxFloat64 {
			return false
		}
	}


	s.matrix.fmap(func(v float64, r, c int) float64 {
		if v != math.MaxFloat64 {
			return v - minForCols[c]
		}
		return v
	})

	s.minBound += sum(minForRows) + sum(minForCols)

	found := false
	s.matrix.each(func(v float64, r int, c int) {
		if v == 0 {
			penalty := s.matrix.getPenalty(r, c)
			if penalty > s.maxPenalty {
				s.maxPenalty = penalty
				s.zr = r
				s.zc = c
				found = true
			}
		}
	})

	return found
}

func (s *Solution) without() *Solution {
	matrix := s.matrix.clone()
	matrix.set(s.zr, s.zc, math.MaxFloat64)
	newSolution := &Solution{
		parent:   s,
		matrix:   matrix,
		minBound: s.minBound + s.maxPenalty,
	}

	return newSolution
}

func (s *Solution) with() *Solution {
	matrix := s.matrix.removeRowAndCol(s.zr, s.zc)
	matrix.set(s.zc, s.zr, math.MaxFloat64)

	newSolution := &Solution{
		parent:   s,
		matrix:   matrix,
		minBound: s.minBound,
		path: append(s.path, pathElem{s.zr, s.zc}),
	}

	return newSolution
}

type PriorityQueue []*Solution

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].minBound > pq[j].minBound
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Solution)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Solution, minBound float64) {
	item.minBound = minBound
	heap.Fix(pq, item.index)
}

func BranchesAndBounds(distances [][]float64, initialWay []int) []int {

	upperBound := math.MaxFloat64

	example := [][]float64{
		{math.MaxFloat64, 26, 42, 15, 29, 25},
		{7, math.MaxFloat64, 16, 1, 30, 25},
		{20, 13, math.MaxFloat64, 35, 5, 0},
		{21, 16, 25, math.MaxFloat64, 18, 18},
		{12, 46, 27, 48, math.MaxFloat64, 5},
		{23, 5, 5, 9, 5, math.MaxFloat64},
	}

	s := &Solution{
		matrix: MatrixFromData(example),
	}

	solutions := PriorityQueue(make([]*Solution, 0))

	heap.Init(&solutions)

	solutions.Push(s)

	for len(solutions) > 0 {
		current := solutions.Pop().(*Solution)
		//fmt.Printf("%s\n\n", current.matrix.debugPrint())
		//fmt.Printf("(%d, %d): %.3f %.3f\n\n", current.zr, current.zc, current.minBound, current.maxPenalty)

		if len(current.matrix.m) == 2 {
			log.Printf("Found full path: %+v", current.path)
			continue
		}

		left := current.without()
		if left.reduce() && left.minBound < upperBound{
			solutions.Push(left)
		}


		right := current.with()
		if right.reduce() && right.minBound < upperBound {
			solutions.Push(right)
		}




	}

	return nil
}

