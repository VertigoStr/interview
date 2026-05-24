// Дан список URL. Нужно выполнить HTTP GET запросы параллельно
// (fan-out), собрать статус-коды (fan-in). 
// Если запрос занимает >1 сек — пропустить.


func FetchStatuses(urls []string) map[string]int {
    in := make(chan string)
    out := make(chan struct{ url string; code int })

    // Fan-out
    for i := 0; i < 5; i++ {
        go func() {
            for url := range in {
                ctx, cancel := context.WithTimeout(context.Background(), time.Second)
                req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
                resp, err := http.DefaultClient.Do(req)
                code := 0
                if err == nil {
                    code = resp.StatusCode
                    resp.Body.Close()
                }
                out <- struct{ url string; code int }{url, code}
                cancel()
            }
        }()
    }

    // Отправка
    go func() {
        for _, url := range urls {
            in <- url
        }
        close(in)
    }()

    // Fan-in
    result := make(map[string]int)
    for i := 0; i < len(urls); i++ {
        r := <-out
        result[r.url] = r.code
    }
    close(out)
    return result
}


