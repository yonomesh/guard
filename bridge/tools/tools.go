package tools

import (
	"context"
	"guard/bridge/tools/x/constraints"
	"runtime"
	"slices"
	"sort"
	"unsafe"
)

// Any checks if any element in the array satisfies the given condition.
func Any[T any](slice []T, fn func(it T) bool) bool {
	for _, it := range slice {
		if fn(it) {
			return true
		}
	}
	return false
}

// AnyIndexed checks if any element in the array satisfies the given condition,
// and it also provides the index of each element to the condition function.
func AnyIndexed[T any](slice []T, fn func(index int, it T) bool) bool {
	for index, it := range slice {
		if fn(index, it) {
			return true
		}
	}
	return false
}

// All checks if all elements in the slice satisfy the given condition.
func All[T any](slice []T, fn func(it T) bool) bool {
	for _, it := range slice {
		if !fn(it) {
			return false
		}
	}
	return true
}

func AllIndexed[T any](slice []T, fn func(index int, it T) bool) bool {
	for index, it := range slice {
		if !fn(index, it) {
			return false
		}
	}
	return true
}

func Contains[T comparable](slice []T, target T) bool {
	return slices.Contains(slice, target)
}

// Map applies the given function fn to each element of the slice
// and returns a new slice of type N with the results.
//
// T: The type of elements in the input slice.
//
// N: The type of elements in the output slice.
//
// Example usage:
//
//	// Define a slice of integers
//	numbers := []int{1, 2, 3, 4}
//
//	// Define a function that doubles an integer
//	double := func(n int) int { return n * 2 }
//
//	// Use Map to double each number in the slice
//	result := Map(numbers, double)
//	fmt.Println(result) // Output: [2, 4, 6, 8]
func Map[T any, N any](slice []T, fn func(it T) N) []N {
	retSlice := make([]N, 0, len(slice))
	for index := range slice {
		retSlice = append(retSlice, fn(slice[index]))
	}

	return retSlice
}

func MapIndexed[T any, N any](slice []T, fn func(index int, it T) N) []N {
	retArr := make([]N, 0, len(slice))
	for index := range slice {
		retArr = append(retArr, fn(index, slice[index]))
	}
	return retArr
}

// FlatMap applies the given mapping function `fn` to each element in the input array `slice`
// and flattens the results into a single array.
func FlatMap[T any, N any](slice []T, fn func(it T) []N) []N {
	var retAddr []N
	for _, item := range slice {
		retAddr = append(retAddr, fn(item)...)
	}
	return retAddr
}

func FlatMapIndexed[T any, N any](slice []T, fn func(index int, it T) []N) []N {
	var retSlice []N
	for index, item := range slice {
		retSlice = append(retSlice, fn(index, item)...)
	}
	return retSlice
}

func Filter[T any](slice []T, fn func(it T) bool) []T {
	var ret []T
	for _, it := range slice {
		if fn(it) {
			ret = append(ret, it)
		}
	}
	return ret
}

func FilterNotNil[T any](slice []T) []T {
	return Filter(slice, func(it T) bool {
		var anyIt any = it
		return anyIt != nil
	})
}

// FilterNotDefault filters out all elements from the slice that are equal to the default value of type T.
//
// Example:
//
//	numbers := []int{1, 0, 3, 0, 5}
//	result := FilterNotDefault(numbers)
//	fmt.Println(result) // output [1 3 5]
func FilterNotDefault[T comparable](arr []T) []T {
	var defaultValue T
	return Filter(arr, func(it T) bool {
		return it != defaultValue
	})
}

// Example:
//
//	arr := []int{10, 20, 30, 40, 50, 60}
//	filtered := FilterIndexed(arr, func(index int, it int) bool {
//		return index%2 == 0 // 保留偶数索引的元素
//	})
//	fmt.Println(filtered) // Output [10 30 50]
func FilterIndexed[T any](slice []T, fn func(index int, it T) bool) []T {
	var retArr []T
	for index, it := range slice {
		if fn(index, it) {
			retArr = append(retArr, it)
		}
	}
	return retArr
}

// The Find function searches for the first element in a slice
// that satisfies a given condition.
// If no even number is found, the function would return the default value (which is 0 for int).
func Find[T any](slice []T, fn func(it T) bool) T {
	for _, it := range slice {
		if fn(it) {
			return it
		}
	}
	return DefaultValue[T]()
}

func FindIndexed[T any](arr []T, fn func(index int, it T) bool) T {
	for index, it := range arr {
		if fn(index, it) {
			return it
		}
	}
	return DefaultValue[T]()
}

func Index[T any](arr []T, fn func(it T) bool) int {
	for index, it := range arr {
		if fn(it) {
			return index
		}
	}
	return -1
}

func IndexIndexed[T any](arr []T, fn func(index int, it T) bool) int {
	for index, it := range arr {
		if fn(index, it) {
			return index
		}
	}
	return -1
}

// The Equal function is used to compare two slices of the same type and determine if they are equal.
// ~[]E It specifies that S must be a slice ([]E) of some element type E, where E is a comparable type.
//
// Example:
//
//	slice1 := []int{1, 2, 3, 4}
//	slice2 := []int{1, 2, 3, 4}
//	slice3 := []int{4, 3, 2, 1}
//
//	fmt.Println(Equal(slice1, slice2)) // Output: true
//	fmt.Println(Equal(slice1, slice3)) // Output: false
func Equal[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// Uniq removes duplicate elements from the input slice and returns a new slice with only unique elements.
func Uniq[T comparable](slice []T) []T {
	result := make([]T, 0, len(slice))
	seen := make(map[T]struct{}, len(slice))

	for _, item := range slice {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// UniqBy removes duplicates from the input slice based on a computed value from each element.
// It uses the provided 'fn' function to generate a comparable value for each element
// and ensures that only unique computed values are kept in the result.
func UniqBy[T any, C comparable](arr []T, fn func(it T) C) []T {
	result := make([]T, 0, len(arr))
	record := make(map[C]struct{}, len(arr))

	for _, item := range arr {
		c := fn(item)
		if _, ok := record[c]; ok {
			continue
		}

		record[c] = struct{}{}

		result = append(result, item)
	}

	return result
}

func SortBy[T any, C constraints.Ordered](arr []T, block func(it T) C) {
	sort.Slice(arr, func(i, j int) bool {
		return block(arr[i]) < block(arr[j])
	})
}

func MinBy[T any, C constraints.Ordered](arr []T, fn func(it T) C) T {
	var min T
	var minValue C
	if len(arr) == 0 {
		return min
	}
	min = arr[0]
	minValue = fn(min)
	for i := 1; i < len(arr); i++ {
		item := arr[i]
		value := fn(item)
		if value < minValue {
			min = item
			minValue = value
		}
	}
	return min
}

func MaxBy[T any, C constraints.Ordered](arr []T, fn func(it T) C) T {
	var max T
	var maxValue C
	if len(arr) == 0 {
		return max
	}
	max = arr[0]
	maxValue = fn(max)
	for i := 1; i < len(arr); i++ {
		item := arr[i]
		value := fn(item)
		if value > maxValue {
			max = item
			maxValue = value
		}
	}
	return max
}

func FilterIsInstance[T any, N any](arr []T, fn func(it T) (N, bool)) []N {
	var retArr []N
	for _, it := range arr {
		if n, isN := fn(it); isN {
			retArr = append(retArr, n)
		}
	}
	return retArr
}

func Reverse[T any](arr []T) []T {
	length := len(arr)
	half := length / 2

	for i := 0; i < half; i = i + 1 {
		j := length - 1 - i
		arr[i], arr[j] = arr[j], arr[i]
	}

	return arr
}

func Done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func Error(_ any, err error) error {
	return err
}

//go:norace
func Dup[T any](obj T) T {
	pointer := uintptr(unsafe.Pointer(&obj))
	//nolint:staticcheck
	//goland:noinspection GoVetUnsafePointer
	return *(*T)(unsafe.Pointer(pointer))

	// 1. unsafe.Pointer(pointer) 将 pointer 变量记录的指针地址转换成 unsafe.Pointer
	// PS unsafe.Pointer 是一个可以指向任意类型的指针
	// 2. (*T)(unsafe.Pointer(pointer)) 将 unsafe.Pointer 类型的指针转换成 T 类型指针
	// 3. *(*T)(unsafe.Pointer(pointer)) 解引用
}

// KeepAlive 的作用是确保传入的对象 obj 在调用 runtime.KeepAlive
// 之后不会被 Go 的垃圾回收器（GC）提前回收。
func KeepAlive(obj any) {
	runtime.KeepAlive(obj)
}

func Must(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}

// Must1 returns result if err is nil, otherwise it panics.
//
// It simplifies handling functions that return (result, error).
//
// Example:
//
//	content := Must1(os.ReadFile("config.json"))
//	fmt.Println(string(content))
func Must1[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}
	return result
}

func Must2[T any, T2 any](result T, result2 T2, err error) (T, T2) {
	if err != nil {
		panic(err)
	}
	return result, result2
}

func PtrOrNil[T any](ptr *T) any {
	if ptr == nil {
		return nil
	}
	return ptr
}

func PtrValueOrDefault[T any](ptr *T) T {
	if ptr == nil {
		return DefaultValue[T]()
	}
	return *ptr
}

func Ptr[T any](obj T) *T {
	return &obj
}

func IsEmpty[T comparable](obj T) bool {
	return obj == DefaultValue[T]()
}

// ClearSlice sets all elements of a Slice to zero value.
func ClearSlice[T ~[]E, E any](t T) {
	clear(t)
}

// ClearMap removes all key-value pairs from the map.
func ClearMap[T ~map[K]V, K comparable, V any](t T) {
	clear(t)
}

// DefaultValue returns the default zero value of type T.
func DefaultValue[T any]() T {
	var defaultValue T
	return defaultValue
}
