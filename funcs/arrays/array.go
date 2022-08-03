package arrays

import "reflect"

// Determine whether a value is included in the array
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i, len := 0, s.Len(); i < len; i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// Get the intersection of two slices
func Intersect(slice1 []int, slice2 []int) []int {
	m := make(map[int]int)
	n := make([]int, 0)
	for _, v := range slice1 {
		m[v]++
	}
	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			n = append(n, v)
		}
	}
	return n
}

// Find out the elements in s1 and not in s2
func Difference(s1, s2 []int) []int {
	m := make(map[int]int)
	n := make([]int, 0)
	inter := Intersect(s1, s2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range s1 {
		if m[value] == 0 {
			n = append(n, value)
		}
	}
	for _, v := range s2 {
		if m[v] == 0 {
			n = append(n, v)
		}
	}
	return n
}

// Remove duplicate values
func RemoveDuplicateWithInt(arr []int) []int {
	var result []int
	tmp := map[int]byte{}
	for _, val := range arr {
		l := len(tmp)
		tmp[val] = 0
		if len(tmp) != l {
			result = append(result, val)
		}
	}
	return result
}

// Remove duplicate values
func RemoveDuplicateWithString(arr []string) []string {
	var result []string
	tmp := map[string]byte{}
	for _, val := range arr {
		l := len(tmp)
		tmp[val] = 0
		if len(tmp) != l {
			result = append(result, val)
		}
	}
	return result
}

// Delete the specified value in the array
func RemoveWithString(arr []string, in string) []string {
	for k, v := range arr {
		if v == in {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	return arr
}

// Delete the specified value in the array
func RemoveWithInt(arr []int, in int) []int {
	for k, v := range arr {
		if v == in {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	return arr
}
