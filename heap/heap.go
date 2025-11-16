package heap

//MaxHeap style -> defo needs optimizations
//Comparator -> returns true if LHS < RHS
func Swap[T any](array []T, index1 int, index2 int) {
	if index1 >= len(array) || index2 >= len(array) {
		return
	}
	temp := array[index1]
	array[index1] = array[index2]
	array[index2] = temp
}

func Heapify[T any](array []T, comparator func(T, T) bool) {
	n := len(array)
	for i := n/2 - 1; i >= 0; i -= 1 {
		Sink(array, i, comparator)
	}
}

func Sink[T any](heap []T, index int, comparator func(T, T) bool) {
	n := len(heap)
	for 2*index+1 < n {
		l := 2*index + 1
		r := 2*index + 2
		cmpPos := l
		if r < n && comparator(heap[l], heap[r]) {
			cmpPos = r
		}
		if comparator(heap[index], heap[cmpPos]) {
			Swap(heap, cmpPos, index)
			index = cmpPos
		} else {
			break
		}
	}

}

func Float[T any](heap []T, index int, comparator func(T, T) bool) {
	for index > 0 {
		cmpPos := index / 2
		if comparator(heap[cmpPos], heap[index]) {
			Swap(heap, cmpPos, index)
			index = cmpPos
		} else {
			break
		}
	}
}

func Push[T any](heap []T, element T, comparator func(T, T) bool) {
	heap = append(heap, element)
	Float(heap, len(heap)-1, comparator)
}

func Pop[T any](heap []T, comparator func(T, T) bool) *T {
	if len(heap) == 0 {
		return nil
	}
	if len(heap) == 1 {
		res := heap[0]
		resPtr := &res
		heap = nil
		return resPtr
	}
	res := heap[0]
	resPtr := &res
	Swap(heap, 0, len(heap)-1)
	heap = heap[:len(heap)-1]
	Sink(heap, 0, comparator)
	return resPtr
}
