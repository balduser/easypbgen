package easypbgen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GenerateTransport(parsed *Parsed) {
	for _, service := range parsed.Services {
		fmt.Printf("Generating service %s\n", service.ServiceName)
		transportFileName := GenerateGRPCTransport(service)
		fmt.Println(transportFileName)
		cmd := exec.Command("gofmt", "-w", transportFileName)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateGRPCTransport(service *Service) string {
	transportFileName := func() string {
		if filename, ok := (*service.Config)[fmt.Sprintf("%sTransportFile", service.ServiceName)]; !ok {
			return fmt.Sprintf("transport%s.go", service.ServiceName)
		} else {
			return filename
		}
	}()
	file, err := os.Create(transportFileName)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(Templates["transportHeading"])
	// Generate adapters
	for _, rpc := range service.RPCs {
		name := rpc.RpcName
		if adapterDeclar, ok := (*service.Config)[fmt.Sprintf(
			"transport%s%s", service.ServiceName, name)]; ok {
			file.WriteString(adapterDeclar)
		} else {
			file.WriteString(fmt.Sprintf(Templates["templateG"], name, name, name, small(name), name, name))
		}
	}
	// Generate decoders and encoders
	for _, rpc := range service.RPCs {
		name := rpc.RpcName
		// Decoders
		if decoderDeclar, ok := (*service.Config)[fmt.Sprintf(
			"transport%sDecode%s", service.ServiceName, rpc.Request.MesName)]; ok {
			file.WriteString(decoderDeclar)
		} else {
			file.WriteString(fmt.Sprintf(
				Templates["templateDec"], name, name, name, genDecFields(&rpc.Request.Fields)))
		}
		// Encoders
		if encoderDeclar, ok := (*service.Config)[fmt.Sprintf(
			"transport%sEncode%s", service.ServiceName, rpc.Response.MesName)]; ok {
			file.WriteString(encoderDeclar)
		} else {
			file.WriteString(fmt.Sprintf(Templates["templateEnc"], name, name, name, genDecFields(&rpc.Request.Fields)))
		}
	}
	return transportFileName
}

// Make first character small
func small(s string) string {
	slice := []rune(s)
	l := strings.ToLower(string(slice[0]))
	return l + s[1:]
}

func genDecFields(fields *[]*Field) string {
	var text string
	for _, field := range *fields {
		text = text + fmt.Sprintf(
			`		%s: req.%s,
`, strings.Title(field.FieldName), strings.Title(field.FieldName))
	}
	return text
}

func genEncFields(fields *[]*Field) string {
	var text string
	for _, field := range *fields {
		switch field.FieldType {
		case "double", "float", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "bool", "string", "bytes":
			text = text + fmt.Sprintf(Templates["encFieldTypeTemplate"], strings.Title(field.FieldName), strings.Title(field.FieldName))
		default:
			text = text + fmt.Sprintf(Templates["encFieldTypeTemplate"], strings.Title(field.FieldName), "XXXXXXXXXXXX")
		}
	}
	return text
}
