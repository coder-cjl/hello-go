package main

import "fmt"

type Logger struct{}

func (l Logger) Debug(a ...any) {
	fmt.Println("[DEBUG]", a)
}

func (l Logger) Info(a ...any) {
	fmt.Println("[INFO]", a)
}

func (l Logger) Error(a ...any) {
	fmt.Println("[ERROR]", a)
}

func (l Logger) Warning(a ...any) {
	fmt.Println("[WARNING]", a)
}

func (l Logger) Fatal(a ...any) {
	fmt.Println("[FATAL]", a)
}

// 提供一个全局实例，可以直接使用 Log.Debug()
var Log = Logger{}
