# httplink [![Go Report Card](https://goreportcard.com/badge/github.com/alexferl/httplink)](https://goreportcard.com/report/github.com/alexferl/httplink) [![codecov](https://codecov.io/gh/alexferl/httplink/branch/master/graph/badge.svg)](https://codecov.io/gh/alexferl/httplink)

A Go module to generate the HTTP [Link](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Link) header.

## Installing
```shell
go get github.com/alexferl/httplink
```

## Using
### Code example
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/alexferl/httplink"
)

func handler(w http.ResponseWriter, r *http.Request) {
	httplink.Append(w.Header(), "/things/2842", "next")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

Make a request:
```shell
curl -i http://localhost:8080
HTTP/1.1 200 OK
Link: </things/2842>; rel=next
Date: Tue, 08 Nov 2022 06:37:22 GMT
Content-Length: 0
```

## Credits
Port of [this](https://github.com/falconry/falcon/blob/3.1.0/falcon/response.py#L779) [Falcon](https://falconframework.org/) method.
