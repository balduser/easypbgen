package easypbgen

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
)

func GenerateTransport(parsed *Parsed) { 
	fmt.Println(parsed.Services)
	for _, service := range parsed.Services {
		fmt.Printf("Generating service %s", service.ServiceName)
		transportFileName := GenerateGRPCTransport(service)
		cmd := exec.Command("gofmt", "-w", transportFileName)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GenerateGRPCTransport(service *Service) string {
	transportFileName := fmt.Sprintf("transport%s.go", service.ServiceName)
	file, err := os.Create(transportFileName)
	if err != nil {
		fmt.Println("Unable to create file:", err) 
        os.Exit(1)
	}
	defer file.Close()
	file.WriteString(initialTemplate)
	// Generate adapters
	for _, rpc := range service.RPCs {
		name := rpc.RpcName
		file.WriteString(fmt.Sprintf(templateG, name, name, name, small(name), name, name))
	}
	// Generate decoders and encoders
	for _, rpc := range service.RPCs {
		file.WriteString(fmt.Sprintf(templateDec, rpc.RpcName, rpc.RpcName, rpc.RpcName, genDecFields(&rpc.Request.Fields) ))
		file.WriteString(fmt.Sprintf(fmt.Sprintf(templateEnc, rpc.RpcName, rpc.RpcName, rpc.RpcName, genEncFields(&rpc.Response.Fields))))
	}
	return transportFileName
}

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
`, strings.Title(field.FieldName), strings.Title(field.FieldName) )
	}
	return text
}

func genEncFields(fields *[]*Field) string {
	var text string
	for _, field := range *fields {
		switch field.FieldType {
        case "double", "float", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "bool", "string", "bytes":
			text = text + fmt.Sprintf(encFieldTypeTemplate, strings.Title(field.FieldName), strings.Title(field.FieldName))
		default:
			text = text + fmt.Sprintf(encFieldTypeTemplate, strings.Title(field.FieldName), "XXXXXXXXXXXX")
		}
	}
	return text
}

const initialTemplate string =
`package transport

import (
	"context"
	// pb
	// model
)

`

const templateG string = 
`func (g grpcTransport) %s(ctx context.Context, request *pb.%sRequest) (*pb.%sResponse, error) {
	_, response, err := g.%s.ServeGRPC(ctx, request)
	if err != nil {
		g.log.Error().Err(err).Msg("%s transport error")
		return nil, err
	}
	resp := response.(*pb.%sResponse)
	return resp, nil
}

`

const templateDec string = 
`func decode%s(ctx context.Context, grpcRequest interface{}) (interface{}, error) {
	req := grpcRequest.(*pb.%sRequest)
	result := &model.%sRequest{
%s	}
	return result, nil
}

`

const templateEnc string = 
`func encode%s(ctx context.Context, grpcResponse interface{}) (interface{}, error) {
	resp := grpcResponse.(*model.%sResponse)
	response := &pb.%sResponse{
%s	}
	return response, nil
}

`

const encFieldTypeTemplate =
`		%s: resp.%s,
`
