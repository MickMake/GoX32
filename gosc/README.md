# GOSC
Gosc is a Go native implementation of the OSC protocol.

Protocol specification can be found at:
[https://ccrma.stanford.edu/groups/osc/spec-1_0.html][1]

## Example

```go
package main

import (
	"github.com/loffa/gosc"
	"log"
)

func main() {
	cli, err := gosc.NewClient("127.0.0.1:1234")
	if err != nil {
		log.Fatalln(err)
	}
	err = cli.EmitMessage("/info")
	if err != nil {
		log.Fatalln(err)
    }
	log.Println("Info message sent!")
}
```

[1]: https://ccrma.stanford.edu/groups/osc/spec-1_0.html