package helper

import "reflect"

// 判断某一个值是否含在切片之中
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

// 取两个切片的交集
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

// 取要校验的和已经校验过的差集，找出需要校验的切片IP（找出slice1中  slice2中没有的）
func Difference(slice1, slice2 []int) []int {
	m := make(map[int]int)
	n := make([]int, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice1 {
		if m[value] == 0 {
			n = append(n, value)
		}
	}

	for _, v := range slice2 {
		if m[v] == 0 {
			n = append(n, v)
		}
	}
	return n
}

//  RemoveDuplicateWithInt 去除重复值
//  @author Cloud|2021-12-12 18:20:52
//  @param arr []int ...
//  @return []int ...
func RemoveDuplicateWithInt(arr []int) []int {
	var result []int      // 存放返回的不重复切片
	tmp := map[int]byte{} // 存放不重复主键
	for _, val := range arr {
		l := len(tmp)
		tmp[val] = 0 // 当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		// 如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tmp) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, val) // 当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

//  RemoveDuplicateWithString 去除重复值
//  @author Cloud|2021-12-13 09:33:05
//  @param arr []string ...
//  @return []string ...
func RemoveDuplicateWithString(arr []string) []string {
	var result []string      // 存放返回的不重复切片
	tmp := map[string]byte{} // 存放不重复主键
	for _, val := range arr {
		l := len(tmp)
		tmp[val] = 0 // 当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		// 如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tmp) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, val) // 当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

//  RemoveWithString 删除数组中的指定值
//  @author Cloud|2021-12-12 18:20:21
//  @param arr []string ...
//  @param in string ...
//  @return []string ...
func RemoveWithString(arr []string, in string) []string {
	for k, v := range arr {
		if v == in {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	return arr
}

//  RemoveWithInt 删除数组中的指定值
//  @author Cloud|2021-12-14 12:15:15
//  @param arr []int ...
//  @param in int ...
//  @return []int ...
func RemoveWithInt(arr []int, in int) []int {
	for k, v := range arr {
		if v == in {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	return arr
}
