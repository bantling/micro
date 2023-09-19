package funcs

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/cmplx"
	"reflect"
	"sort"

	"github.com/bantling/micro/constraint"
	"github.com/bantling/micro/tuple"
)

const (
	notNilableMsg              = "Type %s is not a nillable type"
	sliceFlattenArgNotSliceMsg = "SliceFlatten argument must be a slice, not type %T"
	sliceFlattenArgNotTMsg     = "SliceFlatten argument must be slice of %s, not a slice of %s"
	assertTypeMsg              = "expected %s to be %T, not %T"
	convertToSliceMsg          = "expected %s to be %T, not %T"
	convertToSliceElemMsg      = "expected %s[%d] to be %T, not %T"
	convertToSlice2D1Msg       = "expected %s to be %T, not %T"
	convertToSlice2D2Msg       = "expected %s[%v] to be %T, not %T"
	convertToSlice2ElemMsg     = "expected %s[%v][%v] to be %T, not %T"
	assertMapTypeMsg           = "expected %s to be %T, not %T"
	assertMapTypeValueMsg      = "expected %s[%v] to be %T, not %T"

	// The ALL constant is for some remove funcs
	ALL = true
)

// ==== Slices

// SliceAdd makes adding values to slices easier, by automatically creating the slice as needed.
// Particularly useful for struct fields, as the zero value of the struct will have a nil slice.
//
// EG:
// var slc []int // nil slice
// SliceAdd(&slc, 3)
//
// since slc is nil, sets *slc = []int{}
// appends 3 to *slc
//
// SliceAdd(&slc, 4)
//
// since *slc exists, appends 4 to *slc
func SliceAdd[T any](slc *[]T, value T) {
	if *slc == nil {
		*slc = []T{}
	}

	(*slc) = append((*slc), value)
}

// SliceCopy returns a copy of a slice, useful for situations such as sorting a copy of a slice without modifying the original.
// If the original slice is null, the result is empty.
func SliceCopy[T any](slc []T) (res []T) {
	res = make([]T, len(slc))
	copy(res, slc)
	return
}

// SliceFlatten flattens a slice of any number of dimensions into a one dimensional slice.
// The slice is received as type any, because there is no way to describe a slice of any number of dimensions
// using generics.
//
// If a nil value is passed, an empty slice is returned.
//
// Panics if:
// - the value passed is not a slice
// - the slice passed does not ultimately resolve to elements of type T once all slice dimensions are indexed
func SliceFlatten[T any](value any) []T {
	rslc := []T{}

	// Return empty slice if value is nil
	if value == nil {
		return rslc
	}

	// Make a one dimensional slice to return
	var (
		rtyp = reflect.ValueOf(rslc).Type().Elem()
		vslc = reflect.ValueOf(value)
		vtyp = vslc.Type()
	)

	// Ensure value passed is really a slice
	if vtyp.Kind() != reflect.Slice {
		panic(fmt.Errorf(sliceFlattenArgNotSliceMsg, value))
	}

	// Index all dimensions of value to get the element type
	numDims := 0
	for vtyp.Kind() == reflect.Slice {
		vtyp = vtyp.Elem()
		numDims++
	}

	// Ensure value element type is same as T
	if rtyp != vtyp {
		panic(fmt.Errorf(sliceFlattenArgNotTMsg, rtyp, vtyp))
	}

	// If original value is already one dimenion return it by reference
	if numDims == 1 {
		return value.([]T)
	}

	// Recursively iterate all dimensions of the given slice, some dimensions might be empty
	var f func(reflect.Value)
	f = func(currentSlice reflect.Value) {
		// Iterate current slice
		for i, num := 0, currentSlice.Len(); i < num; i++ {
			val := currentSlice.Index(i)

			// Recurse sub-arrays/slices
			if val.Kind() == reflect.Slice {
				f(val)
			} else {
				rslc = append(rslc, val.Interface().(T))
			}
		}
	}
	f(vslc)

	return rslc
}

// SliceIndex returns the first of the following given a slice, index, and optional default value:
// 1. slice[index] if the slice is non-nil and length > index
// 2. default value if provided
// 3. zero value of slice element type
func SliceIndex[T any](slc []T, index uint, defawlt ...T) T {
	// Return index if it exists
	idx := int(index)
	if (slc != nil) && (len(slc) > idx) {
		return slc[idx]
	}

	// Else return default if provided
	if len(defawlt) > 0 {
		return defawlt[0]
	}

	// Else return zero value
	var zv T
	return zv
}

// Sliceof allows caller to infer the slice type rather than have to write it out.
// This is useful when the type is a more lengthy declaration.
func SliceOf[T any](vals ...T) []T {
	return vals
}

// SliceRemove removes a slice element from a slice.
// By default, only the first occurrence is removed. If the option all param is true, then all occurrences are removed.
// The new slice is returned.
//
// Note: If only the first occurrence is removed, then the append builtin is used twice, once with all elements before
// the occurrence, and again fior all elelemts after it. Otherwise, a new slice is populated with all other elements.
func SliceRemove[T comparable](slc []T, val T, all ...bool) []T {
	// Handle case of first occurrence
	if !SliceIndex(all, 0, false) {
		for i, t := range slc {
			if t == val {
				newSlc := make([]T, 0, len(slc)-1)
				newSlc = append(append(newSlc, (slc)[0:i]...), (slc)[i+1:]...)
				return newSlc
			}
		}

		// No occurrences found, return slc as is
		return slc
	}

	// Handle case of all occurrences
	newSlc := []T{}
	for _, t := range slc {
		if t != val {
			newSlc = append(newSlc, t)
		}
	}

	return newSlc
}

// SliceRemoveUncomparable removes a slice element from a slice of uncomparable values (they must be pointers).
// By default, only the first occurrence is removed. If the optional all param is true, then all occurrences are removed.
// The new slice is returned.
//
// Note: If only the first occurrence is removed, then the the append builtin is used twice, once with all elements before
// the occurrence, and again for all elements after it. Otherwise, a new slice is populated with all other elements.
func SliceRemoveUncomparable[T any](slc []T, val T, all ...bool) []T {
	// Get pointer of val
	valPtr := fmt.Sprintf("%p", any(val))

	// Handle case of first occurrence
	if !SliceIndex(all, 0, false) {
		for i, t := range slc {
			if fmt.Sprintf("%p", any(t)) == valPtr {
				newSlc := make([]T, 0, len(slc)-1)
				newSlc = append(append(newSlc, (slc)[0:i]...), (slc)[i+1:]...)
				return newSlc
			}
		}

		// No occurrences found, return slc as is
		return slc
	}

	// Handle case of all occurrences
	newSlc := []T{}
	for _, t := range slc {
		if fmt.Sprintf("%p", any(t)) != valPtr {
			newSlc = append(newSlc, t)
		}
	}
	return newSlc
}

// SliceReverse reverses the elements of a slice, so that [1,2,3] becomes [3,2,1].
// The slice is modified in place, and returned.
func SliceReverse[T any](slc []T) []T {
	l := len(slc)
	for i := 0; i < l/2; i++ {
		j := l - 1 - i
		tmp := slc[i]
		slc[i] = slc[j]
		slc[j] = tmp
	}

	return slc
}

// SliceSortOrdered sorts a slice of Ordered.
// The slice is modified in place, and returned.
func SliceSortOrdered[T constraint.Ordered](slc []T) []T {
	sort.Slice(slc, func(i, j int) bool { return slc[i] < slc[j] })
	return slc
}

// SliceSortComplex sorts a slice of Complex.
// The slice is modified in place, and returned.
func SliceSortComplex[T constraint.Complex](slc []T) []T {
	sort.Slice(slc, func(i, j int) bool { return cmplx.Abs(complex128(slc[i])) < cmplx.Abs(complex128(slc[j])) })
	return slc
}

// SliceSortCmp sorts a slice of Cmp.
// The slice is modified in place, and returned.
func SliceSortCmp[T constraint.Cmp[T]](slc []T) []T {
	sort.Slice(slc, func(i, j int) bool { return slc[i].Cmp(slc[j]) < 0 })
	return slc
}

// SliceSortBy sorts a slice of any type with the provided comparator.
// The slice is modified in place, and returned.
func SliceSortBy[T any](slc []T, less func(T, T) bool) []T {
	sort.Slice(slc, func(i, j int) bool { return less(slc[i], slc[j]) })
	return slc
}

// SliceToMap copies slice elements to map keys, where the value of each key is true, to avoid the need for two arg map
// lookups to see if a key exists.
// If the slice is nil or empty, an empty map is returned.
func SliceToMap[T comparable](slc []T) (res map[T]bool) {
	res = map[T]bool{}

	for _, val := range slc {
		res[val] = true
	}

	return
}

// SliceToMapBy transforms slice elements to map keys using provided func, where the value of each key is true, to avoid
// the need for two arg map lookups to see if a key exists.
// If the slice is nil or empty, an empty map is returned.
func SliceToMapBy[T any, K comparable](slc []T, fn func(T) K) (res map[K]bool) {
	res = map[K]bool{}

	for _, val := range slc {
		res[fn(val)] = true
	}

	return
}

// SliceUniqueValues returns the uniq values of a slice, in no particular order
// See MapKeysToSlice
func SliceUniqueValues[T comparable](slc []T) []T {
	uniq := map[T]int{}
	for _, v := range slc {
		uniq[v] = 0
	}

	return MapKeysToSlice(uniq)
}

// ==== Maps

// MapIndex returns the first of the following:
// 1. map[key] if the map is non-nil and the key exists in the map
// 2. default if provided
// 3. zero value of map value type
func MapIndex[K comparable, V any](mp map[K]V, key K, defawlt ...V) V {
	// Return key value if it exists
	if mp != nil {
		if val, haveIt := mp[key]; haveIt {
			return val
		}
	}

	// Else return default if provided
	if len(defawlt) > 0 {
		return defawlt[0]
	}

	// Else return zero value of map value type
	var zv V
	return zv
}

// MapSet makes accessing maps easier, by automatically creating the map as needed.
// Particularly useful for struct fields, as the zero value of the struct will have a nil map.
//
// EG:
// var mp map[string]int // nil map
// MapSet(&mp, "foo", 3)
//
// since mp is nil, sets *mp = map[string]{}
// sets *mp["foo"] = 3
//
// MapSet(&mp, "foo", 4)
//
// since *mp exists, sets *mp["foo"] = 4
func MapSet[K comparable, V any](mp *map[K]V, key K, value V) {
	// Create top level map if it is nil
	if *mp == nil {
		*mp = map[K]V{}
	}

	// Set value for key K
	(*mp)[key] = value
}

// MapTest returns true if the map is non-nil and contains the given key
func MapTest[K comparable, V any](mp map[K]V, key K) (exists bool) {
	if mp != nil {
		_, exists = mp[key]
	}

	return
}

// MapUnset is the reverse of MapSet, for consistency
func MapUnset[K comparable, V any](mp map[K]V, key K) {
	// Doesn't matter if map is nil, or if key does not exist
	delete(mp, key)
}

// Map2Set makes accessing two level maps easier, by automatically creating the first and second maps as needed.
// Particularly useful for struct fields, as the zero value of the struct will have a nil map.
//
// EG:
// var mp map[string]map[int]bool // nil map
// Map2Set(&mp, "foo", 3, true)
//
// since mp is nil, sets mp = map[int]bool{}
// since mp["foo"] is nil, sets mp["foo"] = map[int]bool{}
// sets mp["foo"][3] = true
//
// Map2Set(&mp, "foo", 4, false)
//
// since mp and mp["foo"] both exist, sets mp["foo"][4] = false
func Map2Set[K1, K2 comparable, V any](mp *map[K1]map[K2]V, key1 K1, key2 K2, value V) {
	// Create top level map if it is nil
	if *mp == nil {
		*mp = map[K1]map[K2]V{}
	}

	// Create second level map for key K1 if it does not exist
	mp2 := (*mp)[key1]
	if mp2 == nil {
		mp2 = map[K2]V{}
		(*mp)[key1] = mp2
	}

	// Set second level value for key K2
	mp2[key2] = value
}

// Map2Test returns true if mp[K1] and
func Map2Test[K1, K2 comparable, V any](mp map[K1]map[K2]V, key1 K1, key2 K2) (exists bool) {
	if mp != nil {
		if m2 := mp[key1]; m2 != nil {
			_, exists = m2[key2]
		}
	}

	return
}

// Map2Unset makes accessing two level maps easier, by removing K2 from second level map.
//
// EG:
// var mp map[string]map[int]bool
// Map2Set(&mp, "foo", 3, true)
// Map2Set(&mp, "foo", 4, false)
// Map2Unset(&mp, "foo", 3) -> only mp["foo"][4] exists
func Map2Unset[K1, K2 comparable, V any](mp map[K1]map[K2]V, key1 K1, key2 K2) {
	if mp != nil {
		// Doesn't matter if second level map is nil, or if key2 does not exist
		delete(mp[key1], key2)
	}
}

// MapSliceAdd makes accessing a map whose value is a slice easier, by automatically creating the map and slice as needed.
// Particularly useful for struct fields, as the zero value of the struct will have a nil map.
//
// EG:
// var mp map[string][]int
// MapSliceAdd(&mp, "foo", 3)
//
// since mp is nil, sets mp = map[string][]int{}
// since mp["foo"] is nil, sets mp["foo"] = []int{}
// appends 3 to map["foo"]
//
// MapSliceAdd(&mp, "foo", 4)
//
// since mp and mp["foo"] both exist, appends 4 to mp["foo"]
func MapSliceAdd[K comparable, V any](mp *map[K][]V, key K, value V) {
	// Create top level map if it is nil
	if *mp == nil {
		*mp = map[K][]V{}
	}

	// Create slice for key K if it does not exist
	slc := (*mp)[key]
	if slc == nil {
		slc = []V{}
		(*mp)[key] = slc
	}

	// Append value to slice and remap it
	slc = append(slc, value)
	(*mp)[key] = slc
}

// MapSliceRemove makes accessing a map whose value is a slice easier, by removing a value from the slice if it exists.
// If the optional all is true, all occurrences of V are removed, otherwise only the first is removed.
func MapSliceRemove[K, V comparable](mp map[K][]V, key K, value V, all ...bool) {
	if mp != nil {
		if slc, hasIt := mp[key]; hasIt {
			mp[key] = SliceRemove(slc, value, all...)
		}
	}
}

// MapSliceRemoveUncomparable is an uncomparable version of MapSliceRemove
func MapSliceRemoveUncomparable[K comparable, V any](mp map[K][]V, key K, value V, all ...bool) {
	if mp != nil {
		if slc, hasIt := mp[key]; hasIt {
			mp[key] = SliceRemoveUncomparable(slc, value, all...)
		}
	}
}

// Map2SliceAdd makes accessing a two level map whose value is a slice easier, by automatically creating the first and
// second maps and slice as needed.
// Particularly useful for struct fields, as the zero value of the struct will have a nil map.
//
// EG:
// var mp map[string]map[int][]bool
// Map2SliceAdd(&mp, "foo", 3, true)
//
// since mp is nil, sets mp = map[string]map[int][]bool{}
// since mp["foo"] is nil, sets mp["foo"] = map[int][]bool{}
// since mp["foo"][3] is nil, sets mp["foo"][3] = []bool{}
// appends true to map["foo"][3]
//
// Map2SliceAdd(&mp, "foo", 3, false)
//
// since mp, mp["foo"] and mp["foo"][3] all exist, appends false to mp["foo"][3]
func Map2SliceAdd[K1, K2 comparable, V any](mp *map[K1]map[K2][]V, key1 K1, key2 K2, value V) {
	// Create top level map if it is nil
	if *mp == nil {
		*mp = map[K1]map[K2][]V{}
	}

	// Create second level map for key K1 if it does not exist
	mp2 := (*mp)[key1]
	if mp2 == nil {
		mp2 = map[K2][]V{}
		(*mp)[key1] = mp2
	}

	// Create slice for key K2 if it does not exist
	slc := (*mp)[key1][key2]
	if slc == nil {
		slc = []V{}
		(*mp)[key1][key2] = slc
	}

	// Append value to slice and remap it
	slc = append(slc, value)
	(*mp)[key1][key2] = slc
}

// Map2SliceRemove makes accessing a two level map whose value is a slice easier, by removing a value from the slice if it exists.
// If the optional all is true, all occurrences of V are removed, otherwise only the first is removed.
func Map2SliceRemove[K1, K2, V comparable](mp map[K1]map[K2][]V, key1 K1, key2 K2, value V, all ...bool) {
	if mp != nil {
		if mp2, hasIt := mp[key1]; hasIt {
			if slc, hasIt := mp2[key2]; hasIt {
				mp2[key2] = SliceRemove(slc, value, all...)
			}
		}
	}
}

// Map2SliceRemoveUncomparable is an uncomparable version of Map2SliceRemove
func Map2SliceRemoveUncomparable[K1, K2 comparable, V any](mp map[K1]map[K2][]V, key1 K1, key2 K2, value V, all ...bool) {
	if mp != nil {
		if mp2, hasIt := mp[key1]; hasIt {
			if slc, hasIt := mp2[key2]; hasIt {
				mp2[key2] = SliceRemoveUncomparable(slc, value, all...)
			}
		}
	}
}

// MapSortOrdered sorts a map with an Ordered key into a []Two[K, V]
func MapSortOrdered[K constraint.Ordered, V any](mp map[K]V) []tuple.Two[K, V] {
	// Collect pairs into a []Tuple2[K, V]
	var slc []tuple.Two[K, V]
	for k, v := range mp {
		slc = append(slc, tuple.Of2(k, v))
	}

	sort.Slice(slc, func(i, j int) bool { return slc[i].T < slc[j].T })
	return slc
}

// MapSortComplex sorts a map with a Complex key into a []Two[K, V]
func MapSortComplex[K constraint.Complex, V any](mp map[K]V) []tuple.Two[K, V] {
	// Collect pairs into a []Tuple2[K, V]
	var slc []tuple.Two[K, V]
	for k, v := range mp {
		slc = append(slc, tuple.Of2(k, v))
	}

	sort.Slice(slc, func(i, j int) bool { return cmplx.Abs(complex128(slc[i].T)) < cmplx.Abs(complex128(slc[j].T)) })
	return slc
}

// MapSortCmp sorts a map with a Cmp key into a []Two[K, V]
func MapSortCmp[K constraint.Cmp[K], V any](mp map[K]V) []tuple.Two[K, V] {
	// Collect pairs into a []Tuple2[K, V]
	var slc []tuple.Two[K, V]
	for k, v := range mp {
		slc = append(slc, tuple.Of2(k, v))
	}

	sort.Slice(slc, func(i, j int) bool { return slc[i].T.Cmp(slc[j].T) < 0 })
	return slc
}

// MapSortBy sorts a map with any type of key into a []Two[K, V]
func MapSortBy[K comparable, V any](mp map[K]V, less func(K, K) bool) []tuple.Two[K, V] {
	// Collect pairs into a []Tuple2[K, V]
	var slc []tuple.Two[K, V]
	for k, v := range mp {
		slc = append(slc, tuple.Of2(k, v))
	}

	sort.Slice(slc, func(i, j int) bool { return less(slc[i].T, slc[j].T) })
	return slc
}

// MapKeysToSlice collects the map keys into a slice, for scenarios that require a slice of unique values only.
// If the map is empty, an empty slice is returned; the result is never nil.
func MapKeysToSlice[K comparable, V any](mp map[K]V) []K {
	var slc = []K{}
	for k := range mp {
		slc = append(slc, k)
	}

	return slc
}

// ==== []Tuple.Two used as a sorted map

// OrderedTuple2Search searches a []Tuple.Two[K, V] for a Tuple.Two whose K value is the one provided.
// The result is the index of the largest K that is <= the K provided, and true if it is an exact match.
// When false is returned, it means if the given K is to be inserted, it needs to be inserted after the returned index.
//
// # If the slice is nil or empty, then the result is -1, false
//
// An ordered []Tuple.Two[K, V] can be used as a sorted map[K]V
func OrderedTuple2Search[K constraint.Ordered, V any](mp []tuple.Two[K, V], key K) (int, bool) {
	var (
		l, r, m = 0, len(mp) - 1, 0
		k       K
	)

	if r < 0 {
		return -1, false
	}

	// Example of how this works, searching for key 3 in a slice that contains keys 1, 2, 5, 6
	// l (left)   = 0
	// r (right)  = len - 1 = 3
	// m (middle) = 0
	// k = zero value of type K = 0
	//
	// l, r = 0, 3:
	// 0 != 3
	// m = ceil((l + r) / 2) = ceil(1.5) = 2
	// index 2 key = 5, which is > desired key 3, search left of m
	// r = m - 1 = 1
	//
	// l, r = 0, 1:
	// 0 != 1
	// m = ceil((l + r) / 2) = ceil(0.5) = 1
	// index 1 key = 2, which is < desired key 3, search right of m
	// l = m = 1
	//
	// l, r = 1, 1
	// 1 == 1
	// return l, index 1 key (2) == desired key 3
	// return 1, false
	//
	// This indicates add new entry for desired key 3 after index 1
	for l != r {
		// m = ceil((l + r) / 2)
		m = (l + r) / 2
		if (l+r)%2 == 1 {
			m++
		}

		if k = mp[m].T; k > key {
			// too high, search to the left of m
			r = m - 1
		} else {
			// too low, search to the right of m
			l = m
		}
	}

	// if l != r, k is not assigned a new value, so compare original key param to mp[l].T
	return l, key == mp[l].T
}

// OrderedTuple2SliceAdd searches a []Tuple.Two[K, []V] for a tuple with the given key K.
// If such a tuple exists, the value V is added to the end of the slice.
// If no such tuple exists, a new tuple is inserted with a slice containing only the value V provided.
func OrderedTuple2SliceAdd[K constraint.Ordered, V any](mp *[]tuple.Two[K, []V], key K, value V) {
	// Get index of matching tuple, if it exists
	index, matches := OrderedTuple2Search(*mp, key)
	switch {
	case index == -1:
		*mp = []tuple.Two[K, []V]{tuple.Of2(key, []V{value})}

	case matches:
		// index of matching tuple, replace tuple with new tuple that has value appended to slice
		(*mp)[index] = tuple.Of2(key, append((*mp)[index].U, value))

	default:
		// index of tuple that a new tuple should be inserted after
		// due to pointer used for original slice, we have to start building from an empty slice
		*mp = append(append(append([]tuple.Two[K, []V]{}, (*mp)[:index+1]...), tuple.Of2(key, []V{value})), (*mp)[index+1:]...)
	}
}

// OrderedTuple2SliceRemoveKey removes a key from the []Tuple.Two[K, V], if the key exists
func OrderedTuple2SliceRemoveKey[K constraint.Ordered, V any](mp *[]tuple.Two[K, V], key K) {
	if index, matches := OrderedTuple2Search(*mp, key); matches {
		*mp = append((*mp)[:index], (*mp)[index+1:]...)
	}
}

// OrderedTuple2SliceRemoveValue removes a value from the []Tuple.Two[K, []V], if the key exists and the slice contains the value.
// If the optional all is true, then all occurrences of the value are removed, else the first occurrence is removed.
func OrderedTuple2SliceRemoveValue[K constraint.Ordered, V comparable](mp *[]tuple.Two[K, []V], key K, value V, all ...bool) {
	if index, matches := OrderedTuple2Search(*mp, key); matches {
		(*mp)[index] = tuple.Of2(key, SliceRemove((*mp)[index].U, value, all...))
	}
}

// OrderedTuple2SliceRemoveUncomparable removes an uncomparable value from the []Tuple.Two[K, []V], if the key exists and the slice contains the value.
// If the optional all is true, then all occurrences of the value are removed, else the first occurrence is removed.
func OrderedTuple2SliceRemoveUncomparable[K constraint.Ordered, V any](mp *[]tuple.Two[K, []V], key K, value V, all ...bool) {
	if index, matches := OrderedTuple2Search(*mp, key); matches {
		(*mp)[index] = tuple.Of2(key, SliceRemoveUncomparable((*mp)[index].U, value, all...))
	}
}

// ==== Filters - and, not, or

// And converts any number of filter funcs (func(T) bool) into the conjunction of all the funcs.
// Short-circuit logic will return false on the first function that returns false.
// If no filters are provided, the result is a function that always returns true.
func And[T any](filters ...func(T) bool) func(T) bool {
	return func(t T) bool {
		result := true

		for _, nextFilter := range filters {
			if result = nextFilter(t); !result {
				break
			}
		}

		return result
	}
}

// Not (filter func) adapts a filter func (func(T) bool) to the negation of the func.
func Not[T any](filter func(T) bool) func(T) bool {
	return func(t T) bool {
		return !filter(t)
	}
}

// Or converts any number of filter funcs (func(T) bool) into the disjunction of all the funcs.
// Short-circuit logic will return true on the first function that returns true.
// If no filters are provided, the result is a function that always returns true.
func Or[T any](filters ...func(T) bool) func(T) bool {
	return func(t T) bool {
		result := true

		for _, nextFilter := range filters {
			if result = nextFilter(t); result {
				break
			}
		}

		return result
	}
}

// ==== Filters - comparators

// LessThan returns a filter func (func(T) bool) that returns true if it accepts a value that is less than the given value
func LessThan[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t < val
	}
}

// LessThanEqual returns a filter func (func(T) bool) that returns true if it accepts a value that is less than or equal to the given value
func LessThanEqual[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t <= val
	}
}

// Equal returns a filter func (func(T) bool) that returns true if it accepts a value that equals the given value with ==
func Equal[T comparable](val T) func(T) bool {
	return func(t T) bool {
		return t == val
	}
}

// In returns a filter func (func(T) bool) that returns true if it accepts a value that equals any given value with ==
func In[T comparable](val ...T) func(T) bool {
	return func(t T) bool {
		for _, v := range val {
			if t == v {
				return true
			}
		}

		return false
	}
}

// GreaterThan returns a filter func (func(T) bool) that returns true if it accepts a value that is greater than the given value
func GreaterThan[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t > val
	}
}

// GreaterThanEqual returns a filter func (func(T) bool) that returns true if it accepts a value that is greater than or equal to the given value
func GreaterThanEqual[T constraint.Ordered](val T) func(T) bool {
	return func(t T) bool {
		return t >= val
	}
}

// ==== Filters - negative, non-negative, positive

// IsNegative returns a filter func (func(T) bool) that returns true if it accepts a negative value.
func IsNegative[T constraint.Signed]() func(T) bool {
	return func(t T) bool {
		return t < 0
	}
}

// IsNonNegative returns a filter func (func(T) bool) that returns true if it accepts a non-negative value.
func IsNonNegative[T constraint.Signed]() func(T) bool {
	return func(t T) bool {
		return t >= 0
	}
}

// IsPositive returns a filter func (func(T) bool) that returns true if it accepts a positive value.
func IsPositive[T constraint.Signed]() func(T) bool {
	return func(t T) bool {
		return t > 0
	}
}

// ==== Composition

// Compose composes one or more funcs that accept and return the same type into a new function that returns
// f_n(f_n-1( ... (f_1(f_0(x))))). Eg, if three funcs f_0, f_1, f_2 are provided in that order, the resulting
// function returns f_2(f_1(f_0(x))).
func Compose[T any](f0 func(T) T, fns ...func(T) T) func(T) T {
	return func(t T) T {
		res := f0(t)
		for _, fn := range fns {
			res = fn(res)
		}

		return res
	}
}

// Compose2 composes two funcs into a new func that transforms (p -> q -> r)
func Compose2[P, Q, R any](
	f0 func(P) Q,
	f1 func(Q) R,
) func(P) R {
	return func(p P) R {
		return f1(f0(p))
	}
}

// Compose3 composes three funcs into a new func that transforms (p -> ... -> s)
func Compose3[P, Q, R, S any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
) func(P) S {
	return func(p P) S {
		return f2(f1(f0(p)))
	}
}

// Compose4 composes four funcs into a new func that transforms (p -> ... -> t)
func Compose4[P, Q, R, S, T any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
) func(P) T {
	return func(p P) T {
		return f3(f2(f1(f0(p))))
	}
}

// Compose5 composes five funcs into a new func that transforms (p -> ... -> u)
func Compose5[P, Q, R, S, T, U any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
	f4 func(T) U,
) func(P) U {
	return func(p P) U {
		return f4(f3(f2(f1(f0(p)))))
	}
}

// Compose6 composes six funcs into a new func that transforms (p -> ... -> v)
func Compose6[P, Q, R, S, T, U, V any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
	f4 func(T) U,
	f5 func(U) V,
) func(P) V {
	return func(p P) V {
		return f5(f4(f3(f2(f1(f0(p))))))
	}
}

// Compose7 composes seven funcs into a new func that transforms (p -> ... -> w)
func Compose7[P, Q, R, S, T, U, V, W any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
	f4 func(T) U,
	f5 func(U) V,
	f6 func(V) W,
) func(P) W {
	return func(p P) W {
		return f6(f5(f4(f3(f2(f1(f0(p)))))))
	}
}

// Compose8 composes eight funcs into a new func that transforms (p -> ... -> x)
func Compose8[P, Q, R, S, T, U, V, W, X any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
	f4 func(T) U,
	f5 func(U) V,
	f6 func(V) W,
	f7 func(W) X,
) func(P) X {
	return func(p P) X {
		return f7(f6(f5(f4(f3(f2(f1(f0(p))))))))
	}
}

// Compose9 composes nine funcs into a new func that transforms (p -> ... -> y)
func Compose9[P, Q, R, S, T, U, V, W, X, Y any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
	f4 func(T) U,
	f5 func(U) V,
	f6 func(V) W,
	f7 func(W) X,
	f8 func(X) Y,
) func(P) Y {
	return func(p P) Y {
		return f8(f7(f6(f5(f4(f3(f2(f1(f0(p)))))))))
	}
}

// Compose10 composes ten funcs into a new func that transforms (p -> ... -> z)
func Compose10[P, Q, R, S, T, U, V, W, X, Y, Z any](
	f0 func(P) Q,
	f1 func(Q) R,
	f2 func(R) S,
	f3 func(S) T,
	f4 func(T) U,
	f5 func(U) V,
	f6 func(V) W,
	f7 func(W) X,
	f8 func(X) Y,
	f9 func(Y) Z,
) func(P) Z {
	return func(p P) Z {
		return f9(f8(f7(f6(f5(f4(f3(f2(f1(f0(p))))))))))
	}
}

// ==== Ternary

// Ternary returns trueVal if expr is true, else it returns falseVal
func Ternary[T any](expr bool, trueVal T, falseVal T) T {
	if expr {
		return trueVal
	}

	return falseVal
}

// TernaryResult returns trueVal() if expr is true, else it returns falseVal()
func TernaryResult[T any](expr bool, trueVal func() T, falseVal func() T) T {
	if expr {
		return trueVal()
	}

	return falseVal()
}

// ==== Nil

// IsNil generates a filter func (func(T) bool) that returns true if the value given is nil.
// A type constraint cannot be used to describe nillable types at compile time, so reflection is used.
// Panics if T is not a nillable type.
func IsNil[T any]() func(T) bool {
	var n T
	MustBeNillable(reflect.TypeOf(n))

	return func(t T) bool {
		return reflect.ValueOf(t).IsNil()
	}
}

// IsNilValue returns true if the value given is nil
func IsNilValue(val any) bool {
	rv := reflect.ValueOf(val)
	// In the case of an untyped nil any value, reflect.ValueOf() returns Invalid, which means you cannot call IsNil()
	return (!rv.IsValid()) || (Nillable(rv.Type()) && rv.IsNil())
}

// IsNonNil generates a filter func (func(T) bool) that returns true if the value given is non-nil.
// A type constraint cannot be used to describe nillable types at compile time, so reflection is used.
// Panics if T is not a nillable type.
func IsNonNil[T any]() func(T) bool {
	var n T
	MustBeNillable(reflect.TypeOf(n))

	return func(t T) bool {
		return !reflect.ValueOf(t).IsNil()
	}
}

// MustBeNillable panics if Nillable(typ) returns false
func MustBeNillable(typ reflect.Type) {
	if !Nillable(typ) {
		panic(fmt.Errorf(notNilableMsg, typ.Name()))
	}
}

// Nillable returns true if the given reflect.Type represents a chan, func, map, pointer, or slice.
func Nillable(typ reflect.Type) bool {
	nillable := true

	switch typ.Kind() {
	case reflect.Chan:
	case reflect.Func:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
	default:
		nillable = false
	}

	return nillable
}

// ==== Error

// Must panics if the error is non-nil, else returns.
// Useful to wrap calls to functions that return only an error.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustValue panics if the error is non-nil, else returns the value of type T.
// Useful to wrap calls to functions that return a value and an error, where the value is only valid if the error is nil.
func MustValue[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

// MustValue2 panics if the error is non-nil, else returns the values of types T and U.
// Useful to wrap calls to functions that return two values and an error, where the values are only valid if the error is nil.
func MustValue2[T, U any](t T, u U, err error) (T, U) {
	if err != nil {
		panic(err)
	}

	return t, u
}

// MustValue3 panics if the error is non-nil, else returns the values of types T, U, and V.
// Useful to wrap calls to functions that return two values and an error, where the values are only valid if the error is nil.
func MustValue3[T, U, V any](t T, u U, v V, err error) (T, U, V) {
	if err != nil {
		panic(err)
	}

	return t, u, v
}

// ==== Type creation and assertion

// AssertType asserts that v is of type T - intended for cases where T is a scalar type
// returns (v type asserted to T, nil) if v is of type T
// returns (T zero value, error) if v is not of type T
func AssertType[T any](msg string, v any) (T, error) {
	if val, isa := v.(T); isa {
		return val, nil
	}

	var zv T
	return zv, fmt.Errorf(assertTypeMsg, msg, zv, v)
}

// MustAssertType is a must version of AssertType
func MustAssertType[T any](msg string, v any) T {
	return MustValue(AssertType[T](msg, v))
}

// ConvertToSlice asserts that the any value given is a []any, and that all elements are type T, converting to a []T
func ConvertToSlice[T any](msg string, v any) (res []T, err error) {
	var (
		zvElem        T
		slc, isSlcAny = v.([]any)
	)

	if !isSlcAny {
		err = fmt.Errorf(convertToSliceMsg, msg, []any(nil), v)
		return
	}
	tgt := make([]T, len(slc))

	for i, v := range slc {
		val, isa := v.(T)
		if !isa {
			err = fmt.Errorf(convertToSliceElemMsg, msg, i, zvElem, v)
			return
		}

		tgt[i] = val
	}

	res = tgt
	return
}

// MustConvertToSlice is a must version of ConvertToSlice
func MustConvertToSlice[T any](msg string, v any) []T {
	return MustValue(ConvertToSlice[T](msg, v))
}

// ConvertToSlice2 asserts that the any value given is a []any, containing elements of []any, containing elements of T,
// converting to a [][]T
func ConvertToSlice2[T any](msg string, v any) (res [][]T, err error) {
	var (
		zvElem            T
		slcD1, isSlcD1Any = v.([]any)
	)

	if !isSlcD1Any {
		err = fmt.Errorf(convertToSlice2D1Msg, msg, []any(nil), v)
		return
	}
	tgtD1 := make([][]T, len(slcD1))

	for iD1, vD1 := range slcD1 {
		slcD2, isSlcD2Any := vD1.([]any)
		if !isSlcD2Any {
			err = fmt.Errorf(convertToSlice2D2Msg, msg, iD1, []any(nil), vD1)
			return
		}
		tgtD2 := make([]T, len(slcD2))
		tgtD1[iD1] = tgtD2

		for iD2, vD2 := range slcD2 {
			elem, isT := vD2.(T)
			if !isT {
				err = fmt.Errorf(convertToSlice2ElemMsg, msg, iD1, iD2, zvElem, vD2)
				return
			}

			tgtD2[iD2] = elem
		}
	}

	res = tgtD1
	return
}

// MustConvertToSlice2 is a must version of ConvertToSlice2
func MustConvertToSlice2[T any](msg string, v any) [][]T {
	return MustValue(ConvertToSlice2[T](msg, v))
}

// ConvertToMap asserts that the any value given is a map[K]any, and all values of the map are type V, converting to a map[K]V
func ConvertToMap[K comparable, V any](msg string, v any) (res map[K]V, err error) {
	var (
		zvVal         V
		mp, isMapKAny = v.(map[K]any)
	)

	if !isMapKAny {
		err = fmt.Errorf(assertMapTypeMsg, msg, map[K]any(nil), v)
		return
	}
	tgt := map[K]V{}

	for k, v := range mp {
		val, isa := v.(V)
		if !isa {
			err = fmt.Errorf(assertMapTypeValueMsg, msg, k, zvVal, v)
			return
		}

		tgt[k] = val
	}

	res = tgt
	return
}

// MustConvertToMap is a must version of ConvertToMap
func MustConvertToMap[K comparable, V any](msg string, v any) map[K]V {
	return MustValue(ConvertToMap[K, V](msg, v))
}

// ==== Supplier

// SupplierOf generates a func() T that returns the given value every time it is called
func SupplierOf[T any](value T) func() T {
	return func() T {
		return value
	}
}

// SupplierCached generates a func() T that caches the result of the given supplier on the first call.
// Any subseqquent calls return the cached value, guaranteeing the provided supplier is invoked at most once.
func SupplierCached[T any](supplier func() T) func() T {
	var (
		isCached  bool
		cachedVal T
	)

	return func() T {
		if !isCached {
			isCached, cachedVal = true, supplier()
		}

		return cachedVal
	}
}

// ==== FirstValue

// FirstValue2 takes two values and returns only the first one.
// Useful for functions that return two results and you only care about the first one
func FirstValue2[T, U any](t T, u U) T {
	return t
}

// FirstValue3 takes three values and returns only the first one.
// Useful for functions that return three results and you only care about the first one
func FirstValue3[T, U, V any](t T, u U, v V) T {
	return t
}

// ==== SecondValue

// SecondValue2 takes two values and returns only the second one.
// Useful for functions that return two results and you only care about the second one
func SecondValue2[T, U any](t T, u U) U {
	return u
}

// SecondValue3 takes three values and returns only the second one.
// Useful for functions that return three results and you only care about the second one
func SecondValue3[T, U, V any](t T, u U, v V) U {
	return u
}

// ==== ThirdValue

// ThirdValue3 takes three values and returns only the third one.
// Useful for functions that return three results and you only care about the third one
func ThirdValue3[T, U, V any](t T, u U, v V) V {
	return v
}

// ==== TryTo

// IgnoreResult takes a func of no args that returns any type, and generates a func of no args and no return value
// that invokes it.
//
// Useful for TryTo function closers.
func IgnoreResult[T any](fn func() T) func() {
	return func() {
		fn()
	}
}

// TryTo executes tryFn, and if a panic occurs, it executes panicFn.
// If any closers are provided, they are deferred in the provided order before the tryFn, to ensure they get closed even if a panic occurs.
// If any closer returns a non-nil error, any remaining closers are still called, as that is go built in behaviour.
//
// This function simplifies the process of "catching" panics over using reverse order code like the following
// (common in unit tests that want to verify the type of object sent to panic):
//
//	func DoSomeStuff() {
//	  ...
//	  func() {
//	    defer zero or more things that have to be closed before we try to recover from any panic
//	    defer func() {
//	      // Some code that uses recover() to try and deal with a panic
//	    }()
//	    // Some code that may panic, which is handled by above code
//	  }
//	  ...
//	}
func TryTo(tryFn func(), panicFn func(any), closers ...func()) {
	// Defer code that attempts to recover a value - first func deferred is called last, so this func is called after all provided closers
	defer func() {
		if val := recover(); val != nil {
			panicFn(val)
		}
	}()

	// Defer all closers in provided order, so they get called in reverse order as expected
	for _, closerFn := range closers {
		defer closerFn()
	}

	// Execute code that may panic
	tryFn()
}
