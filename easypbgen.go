package easypbgen

import (
	"fmt"
)

type Parsed struct {
	Services []*Service
	Messages map[string]*Message
}

type Service struct {
	ServiceName string
	RPCs        []*Rpc
	MessageList []*Message
}

type Rpc struct {
	RpcName  string
	Request  *Message
	Response *Message
}

type Message struct {
	MesName   string
	Fields    []*Field
	ModelFlag bool
}

type Field struct {
	FieldName string
	FieldType string
	Repeated  bool
	Required  bool
	Optional  bool
	Packed    interface{}
}

type Enum struct {
	EnumName     string
	Constants    map[string]int
	DefaultConst string
	Options      []string
}

/*
// В теории в proto  можно описать несколько services, которые могут пользоваться пересекающимся множеством messages.
// Поскольку выбрана стратегия разделять модель данных на несколько файлов, каждый из которых описывает определённый Service, нужно разделить messages так, чтобы соответствующие им
// структуры GO были описаны в файлах модели, которые относятся к соответствующим services. Для этого метод перебирает Messages всех Rpc службы
// (Service.MessageList формируется при парсинге в rpcInit) и у каждого Message проверяет тип данных поля. Если такой тип есть в parsed.Messages, он добавляется в MessageList
// Это нужно для того, чтобы можно было описать не только модели request и response, но и модели составных типов, если таковые используются.
func (service *Service) FillMessageList(parsed *Parsed) {
	for _, message := range service.MessageList {
		fmt.Printf("Checking message %s", message.MesName)
		for _, field := range message.Fields {
			fmt.Printf("\tfieldName: %s, fieldType: %s\n", field.FieldName, field.FieldType)
			if parsed.Messages[field.FieldType] != nil {
				fmt.Printf("\t\tTrying to append %s type\n", field.FieldType)
				service.AppendToMessageList(parsed.Messages[field.FieldType])
			}
		}
	}
}
*/
func (service *Service) AppendToMessageList(message *Message) {
	//fmt.Printf("Appending to MessageList: service %v, message %v\n", service, message)
	if !service.contains(message) {
		service.MessageList = append(service.MessageList, message)
	}
}

func (service *Service) contains(message *Message) bool {
	for _, m := range service.MessageList {
		if m == message {
			return true
		}
	}
	return false
}

func PrintAll(parsed *Parsed) {
	fmt.Println(parsed.Services)
	for _, service := range parsed.Services {
		fmt.Println(service.ServiceName)
		for _, rpc := range service.RPCs {
			fmt.Printf("\t%s has rpc %s with request %s (%v) and response %s (%v)\n", service.ServiceName, rpc.RpcName, rpc.Request.MesName, &rpc.Request.MesName, rpc.Response.MesName, &rpc.Response.MesName)
		}
	}
	for _, message := range parsed.Messages {
		fmt.Println(message.MesName)
		for _, field := range message.Fields {
			fmt.Printf("%v, ", field)
		}
		fmt.Println()
	}
}
