# easypbgen

An easy code generating tool for transport and model files out from .proto files.

[Русская версия](./README_rus.md)

## Example of usage

1) Create a folder, e.g. pbgen.
2) Put there a .proto file, e.g. service.proto.
3) Create a file (e.g. pbgen.go) with the following content:

```go
package main

import (
  "github.com/balduser/easypbgen"
  "fmt"
  "os"
)

func main() {
  parsed, err := easypbgen.ParseFile(os.Args[1]) 
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  easypbgen.PrintAll(parsed)                          // optional, not required
  easypbgen.GenerateTransport(parsed)
  easypbgen.GenerateModel(parsed)
}
```

4) Move to this folder in CLI.
5) Create a module pbgen:  
`go mod init pbgen`
6) Install easypbgen from github:  
`go mod tidy`
7) Launch code generation with a command  
`go run pbgen.go -service.proto`

Files will appear in the same folder.

## To be done

- Tests!
- Config file
- "reserved" in messages
- "default" in messages
- enum (check, maybe refactor)