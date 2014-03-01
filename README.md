# logstasher

[![wercker status](https://app.wercker.com/status/3e29f62318a057c5d9a48c6b376c5b07/m/ "wercker status")](https://app.wercker.com/project/bykey/3e29f62318a057c5d9a48c6b376c5b07)

logstasher is a Martini middleware that prints logstash-compatiable JSON to an `io.Writer` for each HTTP request.

Here's an example from one of the Go microservices we have at @bikeexchange :

``` json
{
  "@timestamp":"2014-03-01T19:08:06+11:00","@version":1,"method":"GET",
  "path":"/locations/slugs/VIC/Williams-Landing","status":200,"size":238,
  "duration":14.059902000000001,"params":{"country":["au"]}
}
```

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
