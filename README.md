# logstasher

logstasher is a Martini middleware that prints logstash-compatiable JSON to an `io.Writer` for each HTTP request.

Used in conjunction with the [rotating file writer](http://github.com/mipearson/rfw) it allows for rotatable logs ready to feed directly into logstash with no parsing.

### Example

``` go
package main

import (
  "log"

  "github.com/codegangsta/martini"
  "github.com/mipearson/logstasher"
  "github.com/mipearson/rfw"
)

func main() {
  m := martini.Classic()

  logstashLogFile, err := rfw.Open("hello.log", 0644)
  if err != nil {
    log.Fatalln(err)
  }
  defer logstashLogFile.Close()
  m.Use(logstasher.Logger(logstashLogFile))

  m.Get("/", func() string {
    return "Hello world!"
  })
  m.Run()
}
```

```
## logstash.conf
input {
  file {
    path => ["hello.log"]
    codec => "json"
  }
}
```
