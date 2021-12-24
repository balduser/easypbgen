# easypbgen

Простой инструмент для генерации файлов транспорта и модели для фреймворка go-kit из файлов protocol buffers (.proto)

[English version](./README.md)

## Пример использования

1) Создайте папку, к примеру pbgen.
2) Положите в неё .proto-файл, например service.proto.
3) Создайте файл (например, pbgen.go) со следующим содержимым:

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
  easypbgen.PrintAll(parsed)               // необязательно
  easypbgen.GenerateTransport(parsed)
  easypbgen.GenerateModel(parsed)
}
```

4) Перейдите в консоли в эту папку.
5) Создайте модуль pbgen:  
`go mod init pbgen`
6) Установите easypbgen с github:  
`go mod tidy`
7) Запустите генерацию файлов транспорта и модели командой  
`go run pbgen.go -service.proto`

Файлы появятся в этой-же папке.

## Файл конфигурации

Можно создать файл конфигурации, в котором указаны необходимые имена и папки размещения генерируемых файлов, а также прочие предпочтения.
