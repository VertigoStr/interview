1. Fan-Out обработка слайса
Реализуйте функцию, которая принимает слайс целых чисел и количество воркеров. 
Каждый воркер обрабатывает свою непрерывную часть слайса 
(например, умножает на 2) 
и отправляет результаты в канал. 
Главная горутина собирает все результаты и возвращает итоговый слайс.


func processSlice(nums []int, workers int) []int {
    if workers <= 0 {
        workers = 1
    }
    chunkSize := (len(nums) + workers - 1) / workers
    out := make(chan int, len(nums))
    var wg sync.WaitGroup

    for i := 0; i < workers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(nums) {
            end = len(nums)
        }
        wg.Add(1)
        go func(s, e int) {
            defer wg.Done()
            for _, v := range nums[s:e] {
                out <- v * 2
            }
        }(start, end)
    }

    go func() {
        wg.Wait()
        close(out)
    }()

    var res []int
    for v := range out {
        res = append(res, v)
    }
    return res
}

