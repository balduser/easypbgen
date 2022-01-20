package easypbgen

import (
	"bufio"
	//"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseFile(filename string, config *map[string]string) (*Parsed, error) {
	var parsed Parsed
	var mode string
	var modelFlag bool
	var targetService *Service
	var targetMessage *Message
	var targetEnum *Enum
	parsed.Messages = make(map[string]*Message)
	parsed.Enums = make(map[string]*Enum)
	if config != nil {
		parsed.Config = *config
	}
	//Templates = make(map[string]string)
	loadTemplates(&parsed)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
		return nil, err
	}
	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := strings.TrimLeft(fileScanner.Text(), " ")
		// Actions by prefix
		if strings.HasPrefix(line, "}") {
			if mode == "enum" {
				mode = "message"
			} else {
				mode = ""
				targetMessage = nil
			}
		} else if strings.HasPrefix(line, "service") {
			targetService = &Service{ServiceName: strings.Split(line, " ")[1]}
			targetService.Config = &parsed.Config
			parsed.Services = append(parsed.Services, targetService)
		} else if strings.HasPrefix(line, "rpc") {
			rpcInit(line, targetService, parsed.Messages)
		} else if strings.HasPrefix(line, "message") {
			mode = "message"
			targetMessage = messagePtr(strings.Split(line, " ")[1], parsed.Messages, targetService, modelFlag)
		} else if strings.HasPrefix(line, "enum") {
			targetEnum = enumPtr(strings.Split(line, " ")[1], parsed.Enums)
			mode = "enum"
		} else if strings.HasPrefix(line, "/*#") {
			modelFlag = true
		} else if strings.HasPrefix(line, "#*/") {
			modelFlag = false

			// Actions by mode
		} else if mode == "message" {
			fieldFill(line, targetMessage)
		} else if mode == "enum" {
			enumFill(line, targetEnum)
		}
	}
	return &parsed, nil
}

// 1) rpcInit парсит строку; 2) создаёт 2 Message, если таких ещё нет (либо возвращает указатели на существующие);
// 3) создаёт структуру Rpc; 4) добавляет Rpc в targetService.RPCs; 5) добавляет новые messages в targetService.MessageList
func rpcInit(line string, targetService *Service, messages map[string]*Message) {
	re, _ := regexp.Compile(`\([A-Za-z]+\)`)
	name := strings.Split(strings.TrimLeft(line, "rpc "), "(")[0]
	reqResp := re.FindAllString(line, -1)
	req := messagePtr(strings.Trim(reqResp[0], "()"), messages, targetService, false)
	resp := messagePtr(strings.Trim(reqResp[1], "()"), messages, targetService, false)
	rpc := &Rpc{
		RpcName:  name,
		Request:  req,
		Response: resp,
	}
	targetService.RPCs = append(targetService.RPCs, rpc)
	targetService.AppendToMessageList(req)
	targetService.AppendToMessageList(resp)
}

// messagePtr 1) проверяет в Parsed.Messages, существует ли message с таким названием;
//2a) Если да, возвращает указатель на этот message; 2b) Если нет, создаёт его и возвращает указатель на созданную структуру, а также добавляет message в MessageList
func messagePtr(name string, messages map[string]*Message, targetService *Service, modelFlag bool) *Message {
	if messages[name] == nil {
		messages[name] = &Message{
			MesName:   name,
			ModelFlag: modelFlag,
		}
		targetService.AppendToMessageList(messages[name])
	}
	return messages[name]
}

// fieldFill парсит строку описания сообщения и определяет его название, тип данных, вспомогательные инструкции и аннотации
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
	targetMessage.Fields = append(targetMessage.Fields, field)
}

func enumPtr(name string, enums map[string]*Enum) *Enum {
	if enums[name] == nil {
		enums[name] = &Enum{EnumName: name}
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
