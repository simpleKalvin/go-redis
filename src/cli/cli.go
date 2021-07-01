package cli

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type UserInput struct {
	Data    sync.Map
}

var userInput UserInput

var Wg sync.WaitGroup

func init() {
	userInput = UserInput{}
	Wg.Add(1)
	go userInput.CommandStart(&Wg)
	Wg.Wait() // 阻塞主进程不退出
	fmt.Println("all process done...")
}

func (this *UserInput) CommandStart(wg *sync.WaitGroup) {
	defer wg.Done()
	var command, key, value string
	for {
		fmt.Print("root@127.0.0.1 > ")
		fmt.Scanln(&command, &key, &value)
		command = strings.ToUpper(command)
		key     = strings.ToUpper(key)
		value   = strings.ToUpper(value)
		outputData, err := this.SwitchCommand(command, key, value)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 返回值
		fmt.Println(outputData)
	}
}

func (this *UserInput) SwitchCommand(command, key, value string) (data interface{}, err error) {
	switch command {
	case "GET":
		// 获取数据
		redisData, err := this.GetData(key)
		return redisData, err
	case "SET":
		redisData, err := this.SetData(key, value)
		return redisData, err
	case "EXIT":
		// 退出
		return nil, errors.New("cli exits")
	default:
		// 退出
		fmt.Println("异常输入.. 请重新输入...")
		return nil, nil
	}
}

func (this *UserInput) GetData(key string) (data interface{}, err error)  {
	load, ok := this.Data.Load(key)
	if !ok {
		return nil, nil
	}
	return load, nil
}


func (this *UserInput) SetData(key, value string) (data interface{}, err error) {
	if key != "" && value != "" {
		this.Data.Store(key, value)
		return "success...", nil
	}
	return "please input values", nil
}