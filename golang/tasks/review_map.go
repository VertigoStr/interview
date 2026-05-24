// Найти проблемы в коде

package main

import (
    "fmt"
    "sync"
)

func main() {
    fmt.Println(GetOrCreate("hello", "world")) //20%
    fmt.Println(Get("hello")) //80%
}

var cache = make(map[string]string)

// GetOrCreate проверяет существование ключа key
// Если такого нет, то создает новое значение
func GetOrCreate(key, value string) string {
    var m sync.Mutex
    m.Lock()
    value = cache[key]
    m.Unlock()
    if value != "" {
        return value
    }
    m.Lock()
    cache[key] = value
    m.Unlock()
    return value
}
func Get(key string) string {
    var m sync.Mutex
    m.Lock()
    v := cache[key]
    m.Unlock()
    return v
}
