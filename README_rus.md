# easypbgen

Простой инструмент для генерации файлов транспорта и модели для фреймворка go-kit из файлов protocol buffers (.proto)

[English version](./README.md)

---
## Пример использования

1) Создать папку, к примеру pbgen
2) Положить в неё .proto-файл, например service.proto
3) Создать в этой папке файл (например, pbgen.go) со следующим содержимым:

```go
package main

import (
  "github.com/balduser/easypbgen"
  "fmt"
  "os"
)

func main() {
  pbgen, err := easypbgen.ParseFile(os.Args[1], nil)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  easypbgen.PrintAll(pbgen)          // необязательно
  easypbgen.GenerateTransport(pbgen)
  easypbgen.GenerateModel(pbgen)
}
```

4) Перейти в консоли в эту папку.
5) Создать модуль pbgen:  
`go mod init pbgen`
6) Установить easypbgen с github:  
`go mod tidy`
7) Запустить генерацию файлов транспорта и модели командой  
`go run pbgen.go -service.proto`

Файлы появятся в этой-же папке.

---

## На заметку

- Нужно, чтобы в proto-файле была как минимум одна Service
- Service должны быть определены до Messages

---

## Конфигурация и шаблоны

Невозможно создать инструмент под все случаи жизни, поэтому код, сгенерированный easypbgen может не быть оптимальным. Некоторые части сгенерированного кода транспорта и модели могут не соответствовать требованиям проекта. Поэтому чтобы иметь возможность автоматической генерации структур в нужном виде, запускать easypbgen многократно и после каждого запуска не править сгененрированный код вручную, можно размещать части требуемого на выходе кода в переменных вызывающего скрипта, или выностить их в конфигурационный файл.  
Для размещения параметров в теле скрипта нужно создать map[string]string, указатель на который нужно передать вторым аргументом в функцию ParseFile(). Создавать можно следующие ключи:

### Параметры выходных файлов

- `<service_name>TransportFile` - uri генерируемого файла кода транспорта

- `<service_name>ModelFile` - uri генерируемого файла кода модели

### Указание частей кода вручную

- `model<service_name>Message<message_name>` - код определённого сообщения

- `transport<service_name><method_name>` - код определённого метода транспортного уровня (адаптера)

- `transport<service_name>Decode<message_name>` - код определённого декодера

- `transport<service_name>Encode<message_name>` - код определённого энкодера

- `model<service_name>Ending` - текст или код для вставки в конец файла модели

### Переопределение встроенных шаблонов

- `transportHeading` - заголовок генерируемого файла транспорта
- `templateG` - шаблон адаптера транспортного уровня
- `templateDec` - шаблон декодера транспортного уровня
- `templateEnc` - шаблон энкодера транспортного уровня
- `encFieldTypeTemplate` - шаблон поля "тип" для энкодера
- `modelHeading` - заголовок генерируемого файла модели
- `modelTemplate` - шаблон структуры в файле модели

<details><summary>Примеры</summary>

---

- код определённого сообщения

```go
"modelBlogpostAPIServiceMessageCreatePostRequest":
`type CreatePostRequest struct {
  userID         uint32
  postText       string
  hellobug       float64
}

`
```

- код определённого метода транспортного уровня (адаптера)

```go
func (g grpcTransport) CreatePost(ctx context.Context, request *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
  _, response, err := g.createPost.ServeGRPC(ctx, request)
  if err != nil {
    g.log.Error().Err(err).Msg("WE LOVE BUGS!")
    return nil, err
  }
  resp := response.(*pb.CreatePostResponse)
  return resp, nil
}

```

- Код определённого декодера

```go
func decodeCreatePost(ctx context.Context, grpcRequest interface{}) (interface{}, error) {
  req := grpcRequest.(*pb.CreatePostRequest)
  result := &model.CreatePostRequest{
    UserID:         req.UserID,
    PostText:       req.PostText,
    hellobug:       float64,
  }
  return result, nil
}
```

- Код определённого энкодера

```go
func encodeCreatePost(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
  resp := grpcResponse.(*model.CreatePostResponse)
  response := &pb.CreatePostResponse{
    CreatePostSuccess: resp.CreatePostSuccess,
    ErrorText:         resp.ErrorText,
    hellobug:          0,
  }
  return response, nil
}

```

---
</details>

Значения шаблонов по умолчанию находятся в файле templates.go

<details><summary>Пример скрипта с параметрами конфигурации в теле</summary>

---

```go
package main

import (
  "github.com/balduser/easypbgen"
  "fmt"
  "os"
)

func main() {
  config := map[string]string {
    "transportHeading": `package transport
// Hello there!
`,
    "modelBlogpostAPIServiceEnding":
`type CreatePostRequest struct {
    myPrettyFieldName    bug64
}
`,
  }

  pbgen, err := easypbgen.ParseFile(os.Args[1], &config)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  easypbgen.GenerateTransport(pbgen)
  easypbgen.GenerateModel(pbgen)
}

```

---

</details>

### Вынос шаблонов кода в файл конфигурации

Для повышения читаемости запускающего скрипта шаблоны кода можно вынести в файл конфигурации, который представляет собой файл .go, размещённый в папке скрипта. Он должен также относиться к package main, и в нём нужно создать map[string]string, указатель на который будет передан в ParseFile(). В таком случае в вызывающем скрипте создавать map не нужно.
В случае использования файла конфигурации проект перед запуском необходимо собрать при помощи `go build`.

---
## Генерация кода модели

Бывают ситуации, когда модели должна охватывать больше структур, чем код транспорта. На этот случай в .proto-файле можно размещать комментарии для модели данных. Они будут учтены функцией ParseFile(), но не повлияют на компиляцию grpc файлов при помощи protoc. Комментарии для модели создаются при помощи сочетания символов `/*#` и `#*/`, обязательно отдельной строкой:

```proto3
...
// Код protocol buffers

/*#
message LikePostRequest {
  uint32 postID = 1;
}

message LikeCommentRequest {
  uint32 commentID = 1;
}

message SuccessResponse {
  bool Success = 1; // true for success, false for failure
  string errorText = 2;
}
#*/
```

Если же вы не хотите раскрывать детали о дополнительных полях модели в proto-файле, дополнительный код таких полей можно разместить в [файле конфигурации](#конфигурация-и-шаблоны) с ключом `model<service_name>Ending`. Это подставит код полей в конец файла.
