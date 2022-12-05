## Introduction

chiab is a simple stress testing tool similar to ab, written in go. ab is a standalone command-line tool, and chiab is a function that is executed embedded in your code. 

## Run()

parameters

* **handler** - code snippet
* **concurrency** - Number of multiple requests to make at a time
* **requests** - Number of requests to perform
* **title** - Reporting title
* **save** - If true save reporting file

return

Returns true on success

## example

testing http get request

```go
func main() {
    var concurrency int64 = 20
    var requests int64 = 10000
    RequestStart(concurrency, 60*time.Second)
    Run(func(id int64) bool {
        _, err := Get(id, "http://127.0.0.1:8100/profile", nil, "")
        if err != nil {
            return false
        } else {
            return true
        }
    }, concurrency, requests, "test http get", false)
}
```

testing code snippet
```go
func main() {
    var concurrency int64 = 20
    var requests int64 = 10000
    RequestStart(concurrency, 60*time.Second)
    Run(func(id int64) bool {
        return true
    }, concurrency, requests, "测试函数执行效率", true)
}
```