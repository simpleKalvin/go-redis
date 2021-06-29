package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type UserInput struct {
	Command string
	Key     string
	Value   string
	Expire  int
	Data    sync.Map
	Ch      chan bool
}

var (
	bgPath = "./data/data.data"
	flagChan chan bool
)

func NewRedis() *UserInput {
	return loadBackUp()
}

func loadBackUp() *UserInput {
	userInput := &UserInput{}
	file, err := os.Open(bgPath)
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("service init failed:", err)
	}

	toRestoreData := make(map[string]string, 10)
	json.Unmarshal(fileData, &toRestoreData)
	for key, value := range toRestoreData {
		userInput.Data.Store(key, value)
	}
	return userInput
}

func (this *UserInput) GET() {

	inputReader := bufio.NewReader(os.Stdin) //创建一个读取器，并将其与标准输入绑定。
	var commandList []string
	for {
		output := "redis 127.0.0.1:6379 > "
		fmt.Print(output)
		//successNum, err := fmt.Scanf("%v %v %v %v", &this.Command, &this.Key, &this.Value, &this.Expire)
		inputString, err := inputReader.ReadString('\n') //读取器对象提供一个方法 ReadString(delim byte) ，该方法从输入中读取内容，直到碰到 delim 指定的字符，然后将读取到的内容连同 delim 字符一起放到缓冲区。
		if err == nil && inputString == ""{
			fmt.Printf("The input was: %s", inputString)
		}
		inputString = strings.TrimRight(inputString, "\n")
		commandList = strings.Split(inputString, " ")
		this.transUserInputCommand(commandList)
	}
}

func (this *UserInput) transUserInputCommand(commandList []string) {
	if commandList[0] == "" {
		return
	}

	command := strings.ToUpper(commandList[0])
	switch command {
	case "GET":
		load, ok := this.Data.Load(commandList[1])
		if ok {
			fmt.Println(load)
			return
		}
		fmt.Println(nil)
	case "SET":
		fmt.Println("ok")
		this.Data.Store(commandList[1], commandList[2])
	case "SAVE":
		fallthrough
	case "BGSAVE":
		go this.backUp()
	default:
		fmt.Println("unknow command")
	}

	this.Command = command
}

func (this *UserInput)backUp() {
	fmt.Println("备份")
	backupData := make(map[string]string, 10)
	this.Data.Range(func(key, value interface{}) bool {
		backupData[key.(string)] = value.(string)
		return true
	})
	backUpString, err := json.Marshal(backupData)
	if err !=nil {
		log.Fatalf("failed: %v", err)
	}

	ioutil.WriteFile(bgPath, []byte(backUpString),  0777)
}
