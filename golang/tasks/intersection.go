

Напишите функцию Intersection(a, b []int) []int, 
которая возвращает слайс из элементов, 
присутствующих в обоих слайсах 
(порядок не важен, дубликаты убрать).


func main() {

	a := []int{23, 3, 1, 2}
	b := []int{6, 2, 4, 23}
	// [2, 23]
	fmt.Printf("%v\n", Intersection(a, b))
	a = []int{1, 1, 1}
	b = []int{1, 1, 1, 1}
	// [1]
	fmt.Printf("%v\n", Intersection(a, b))
}


func Intersection(a, b []int) []int {
    set := make(map[int]struct{})
    res := make([]int, 0)

    for _, v := range a {
        set[v] = struct{}{}
    }

    for _, v := range b {
        if _, ok := set[v]; ok {
            res = append(res, v)
            delete(set, v) // убираем повторное добавление
        }
    }
    return res
}

