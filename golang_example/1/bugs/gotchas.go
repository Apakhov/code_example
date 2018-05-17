package main

import (
	"sort"
	"strconv"
)

func ReturnInt() int {
	return 1
}

func ReturnFloat() float32 {
	return 1.1
}

func ReturnIntArray() [3]int {
	return [3]int{1, 3, 4}
}

func ReturnIntSlice() []int {
	return []int{1, 2, 3}
}

func IntSliceToString(slice []int) string {
	s := ""
	for _, x := range slice {
		s += strconv.Itoa(x)
	}
	return s
}

func MergeSlices(fslice []float32, islice []int32) []int {
	ptf := 0
	pti := 0
	lenf := len(fslice)
	leni := len(islice)
	nslice := make([]int, lenf+leni)

	for i := range nslice {
		switch {
		case (ptf == lenf):
			nslice[i] = int(islice[pti])
			pti++
		case (pti == leni):
			nslice[i] = int(fslice[ptf])
			ptf++
		default:
			if int(fslice[ptf]) < int(islice[pti]) {
				nslice[i] = int(fslice[ptf])
				ptf++
			} else {
				nslice[i] = int(islice[pti])
				pti++
			}
		}
	}
	return nslice
}

type byVal []int

func (bv byVal) Len() int {
	return len(bv)
}
func (bv byVal) Less(i, j int) bool {
	return bv[i] < bv[j]
}
func (bv byVal) Swap(i, j int) {
	bv[i], bv[j] = bv[j], bv[i]
}

func GetMapValuesSortedByKey(mp map[int]string) []string {
	keys := make([]int, 0)
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Sort(byVal(keys))
	answ := make([]string, 0)
	for _, key := range keys {
		answ = append(answ, mp[key])
	}
	return answ
}

// сюда вам надо писать функции, которых не хватает, чтобы проходили тесты в gotchas_test.go
