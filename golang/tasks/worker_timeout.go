Каждая задача принимает context.Context и выполняет 
потенциально долгую операцию. 
Реализуйте пул, который запускает задачи 
с ограничением по времени (передаётся параметром).
Если задача не уложилась – возвращается context.DeadlineExceeded,
остальные задачи продолжаются. Верните слайс ошибок.



type task struct {
    fn func(context.Context) error
}

func runWithTimeout(tasks []task, workers int, timeout time.Duration) []error {
    in := make(chan task)
    res := make(chan error, len(tasks))

    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for t := range in {
                ctx, cancel := context.WithTimeout(context.Background(), timeout)
                err := t.fn(ctx)
                cancel()
                res <- err
            }
        }()
    }

    for _, t := range tasks {
        in <- t
    }
    close(in)
    wg.Wait()
    close(res)

    var errs []error
    for e := range res {
        if e != nil {
            errs = append(errs, e)
        }
    }
    return errs
}
