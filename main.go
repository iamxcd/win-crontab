package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/exec"
	"time"
)

/**
任务
*/
type Task struct {
	Name string
	Cron string
	Cmd  string
}

var logger *log.Logger
var logFile *os.File

func (task Task) Run() {
	cmd := exec.Command("cmd")
	in := bytes.NewBuffer(nil)
	out := bytes.NewBuffer(nil)
	cmd.Stdin = in   //绑定输入
	cmd.Stdout = out //绑定输出

	in.WriteString(task.Cmd + "\n")

	if err := cmd.Run(); err != nil {
		logger.Println("启动命令行错误:", err.Error())
	}
	rt := mahonia.NewDecoder("gbk").ConvertString(out.String()) // cmd默认是gbk 转成utf-8
	logger.Println(time.Now().Format("2006-01-02:15:04:05"))
	logger.Println(rt)
}
func main() {
	tasks, err := readTask()
	if err != nil {
		logger.Println("读取配置文件错误:", err.Error())
	}
	// 初始化日志文件
	initLog()

	cron := cron.New(cron.WithSeconds())
	for _, task := range tasks {
		fmt.Println(task)
		if _, err := cron.AddJob(task.Cron, task); err != nil {
			logger.Println("添加计划任务失败:", err.Error())
		}
	}

	cron.AddFunc("0 0 * * * *", func() {
		logFile.Close()
		initLog()
	})
	fmt.Println("任务执行中,请勿关闭窗口")
	cron.Run()
}

func readTask() ([]Task, error) {

	var tasks []Task
	file, err := os.Open("./cron.json")
	if err != nil {
		return tasks, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	return tasks, nil
}

func initLog() {

	if _, err := os.Stat("logs"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("logs", 0777)
		}
	}

	var err error
	filename := time.Now().Format("2006-01-02~15") + ".log"

	logFile, err = os.Create("./logs/" + filename)

	if err != nil {
		fmt.Println("日志文件创建失败" + err.Error())
		os.Exit(0)
	}

	logger = log.New(logFile, "[Debug]", log.Lshortfile)
}
