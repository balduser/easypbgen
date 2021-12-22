package easypbgen

import (
	"fmt"
)

func GenerateGRPCModel(service *Service) string {
	modelFileName := fmt.Sprintf("model%s.go", service.ServiceName)
	
	return modelFileName
}

func GenerateModel(services *[]*Service) {
	for _, service := range *services {
		fmt.Printf("Generating model %s", service.ServiceName)
		//modelFileName := GenerateGRPCModel(service)
	}
}
