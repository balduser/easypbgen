package easypbgen

import (
	"fmt"
)

type Parsed struct {
	Services     []*Service
	Messages     map[string]*Message
}

type Service struct {
	ServiceName  string
	RPCs         []*Rpc
}

type Rpc struct {
	RpcName      string
	Request      *Message
	Response     *Message
	Messages     []*Message
}

type Message struct {
	MesName      string
	Fields       []*Field
}

type Field struct {
	FieldName    string
	FieldType    string
	Repeated     bool
	Required     bool
	Optional     bool
	Packed       interface{}
}

type Enum struct {
	EnumName     string
	Constants    map[string]int
	DefaultConst string
	Options      []string
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
