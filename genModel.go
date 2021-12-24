package easypbgen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func GenerateModel(parsed *Parsed) {
	for _, service := range parsed.Services {
		fmt.Printf("Generating model for service %s", service.ServiceName)
		modelFileName := GenerateGRPCModel(service)
		cmd := exec.Command("gofmt", "-w", modelFileName)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateGRPCModel(service *Service) string {
	modelFileName := fmt.Sprintf("model%s.go", service.ServiceName)
	file, err := os.Create(modelFileName)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(modelHeading)
	for _, message := range service.MessageList {
		file.WriteString(fmt.Sprintf(modelTemplate, message.MesName, genModelFields(message)))
	}
	return modelFileName
}

func genModelFields(message *Message) string {
	answer := ""
	for _, field := range message.Fields {
		answer = answer + fmt.Sprintf(
			`		%s %s
`, field.FieldName, genModelFieldType(field))
	}
	return answer
}

func genModelFieldType(field *Field) string {
	prefix := ""
	if field.Repeated {
		prefix = "[]"
	}
	fieldType := ""
	switch field.FieldType { // https://developers.google.com/protocol-buffers/docs/overview#scalar
	case "double":
		fieldType = "float64"
	case "float":
		fieldType = "float32"
	case "int32", "sint32", "sfixed32":
		fieldType = "int32"
	case "int64", "sint64", "sfixed64":
		fieldType = "int64"
	case "uint32", "fixed32":
		fieldType = "uint32"
	case "uint64", "fixed64":
		fieldType = "uint64"
	case "bytes":
		fieldType = "[]byte"
	default:
		fieldType = field.FieldType
	}
	answer := prefix + fieldType
	return answer
}

const modelHeading string = `package model
//Generated by github.com/balduser/easypbgen

`

const modelTemplate string = `type %s struct {
%s
}

`
