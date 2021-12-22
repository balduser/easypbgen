package easypbgen

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseFile(filename string) (*Parsed, error) { 
	var parsed        Parsed
	var mode          string
	var targetService *Service
	var targetMessage *Message
	var targetEnum    *Enum
	parsed.Messages = make(map[string]*Message)
	Enums := make(map[string]*Enum)

	file, err := os.Open(filename)
	defer file.Close() 
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		return nil, err
    }
	fileScanner := bufio.NewScanner(file)
	
	for fileScanner.Scan() {
		line := strings.TrimLeft(fileScanner.Text(), " ")
		if strings.HasPrefix(line, "}") {
			if mode == "enum" {
				mode = "message"
			} else {
				mode = ""
				//targetService = nil
				targetMessage = nil
			}
		} else if strings.HasPrefix(line, "service") {
			targetService = &Service{ ServiceName: strings.Split(line, " ")[1] }
			parsed.Services = append(parsed.Services, targetService)
		} else if strings.HasPrefix(line, "rpc") {
			rpcInit(line, targetService, parsed.Messages)
		} else if strings.HasPrefix(line, "message") {
			targetMessage = messageInit(strings.Split(line, " ")[1], parsed.Messages)
			mode = "message"
		} else if strings.HasPrefix(line, "enum") {
			targetEnum = enumInit(strings.Split(line, " ")[1], Enums)
			mode = "enum"
		} else if mode == "message" {
			fieldFill(line, targetMessage)
		} else if mode == "enum" {
			enumFill(line, targetEnum)
		}
	}
	return &parsed, nil
}

func rpcInit(line string, targetService *Service, messages map[string]*Message) {
	re, _ := regexp.Compile(`\([A-Za-z]+\)`)
	name := strings.Split(strings.TrimLeft(line, "rpc "), "(")[0]
	reqResp := re.FindAllString(line, -1)
	req := messageInit(strings.Trim(reqResp[0], "()"), messages)
	resp := messageInit(strings.Trim(reqResp[1], "()"), messages)
	rpc := &Rpc{
		RpcName: name,
		Request: req,
		Response: resp,
	}
	targetService.RPCs = append(targetService.RPCs, rpc)
}

// messageInit 1) проверяет в Messages, существует ли message с таким названием; 2a) Если да, возвращает указатель на этот message; 2b) Если нет, создаёт его и возвращает указатель на созданную структуру  
func messageInit(name string, messages map[string]*Message) *Message {
	if messages[name] == nil {
		messages[name] = &Message{ MesName: name }
	}
	return messages[name]
}

// fieldMaker парсит строку описания сообщения и определяет его название, тип данных, вспомогательные инструкции и аннотации
func fieldFill(line string, targetMessage *Message) {
	field := &Field{}
	if strings.HasPrefix(line, "required") { // https://developers.google.com/protocol-buffers/docs/overview#simple
		field.Required = true
		line = strings.Trim(strings.TrimLeft(line, "required"), " ")
	}
	if strings.HasPrefix(line, "optional") {
		field.Optional = true
		line = strings.Trim(strings.TrimLeft(line, "optional"), " ")
	}
	if strings.HasPrefix(line, "repeated") {
		field.Repeated = true
		line = strings.Trim(strings.TrimLeft(line, "repeated"), " ")
	}
	if strings.Contains(line, `[packed = true]`) { // https://developers.google.com/protocol-buffers/docs/overview#specifying_field_rules
		field.Packed = true
	}
	// To be done: parse "reserved" - https://developers.google.com/protocol-buffers/docs/overview#reserved
	// To be done: parse "default" 
	mType := strings.Split(line, " ")[0]
	mName := strings.Split(line, " ")[1]
	field.FieldName = mName
	field.FieldType = mType
	targetMessage.Fields = append(targetMessage.Fields, field )
}

func enumInit(name string, enums map[string]*Enum) *Enum{
	if enums[name] == nil {
		enums[name] = &Enum{ EnumName: name }
	}
	return enums[name]
}

func enumFill(line string, targetEnum *Enum) { // Not tested!
	if strings.HasPrefix(line, "option") {
		targetEnum.Options = append(targetEnum.Options, strings.TrimLeft(line, "option "))
	}
	name := strings.Split(line, " ")[0]
	number, _ := strconv.Atoi(strings.TrimRight(strings.Split(line, " ")[1], ";"))
	targetEnum.Constants[name] = number
}

