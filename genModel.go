package easypbgen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GenerateModel(parsed *Parsed) {
	for _, service := range parsed.Services {
		fmt.Printf("Generating model for service %s\n", service.ServiceName)
		modelFileName := GenerateGRPCModel(service)
		fmt.Println(modelFileName)
		cmd := exec.Command("gofmt", "-w", modelFileName)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateGRPCModel(service *Service) string {
	modelFileName := func() string {
		if filename, ok := (*service.Config)[fmt.Sprintf("%sModelFile", service.ServiceName)]; !ok {
			return fmt.Sprintf("model%s.go", service.ServiceName)
		} else {
			return filename
		}
	}()
	file, err := os.Create(modelFileName)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(Templates["modelHeading"])
	for _, message := range service.MessageList {
		if messageDeclar, ok := (*service.Config)[fmt.Sprintf("model%sMessage%s", service.ServiceName, message.MesName)]; ok {
			file.WriteString(messageDeclar)
		} else {
			file.WriteString(fmt.Sprintf(Templates["modelTemplate"], message.MesName, genModelFields(message)))
		}
	}
	if modelEnding, ok := (*service.Config)[fmt.Sprintf("model%sEnding", service.ServiceName)]; ok {
		file.WriteString(modelEnding)
	}
	return modelFileName
}

func genModelFields(message *Message) string {
	answer := ""
	for _, field := range message.Fields {
		answer = answer + fmt.Sprintf(
			`		%s %s
`, strings.Title(field.FieldName), genModelFieldType(field))
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
