package main

import "fmt"

type Logger struct{}

func (l Logger) Debug(a ...any) {
	fmt.Println("[DEBUG]", a)
}

func (l Logger) Debugf(format string, a ...any) {
	fmt.Printf("[DEBUG] "+format, a...)
}

func (l Logger) Info(a ...any) {
	fmt.Println("[INFO]", a)
}

func (l Logger) Infof(format string, a ...any) {
	fmt.Printf("[INFO] "+format, a...)
}

func (l Logger) Error(a ...any) {
	fmt.Println("[ERROR]", a)
}

func (l Logger) Errorf(format string, a ...any) {
	fmt.Printf("[ERROR] "+format, a...)
}

func (l Logger) Warning(a ...any) {
	fmt.Println("[WARNING]", a)
}

func (l Logger) Warningf(format string, a ...any) {
	fmt.Printf("[WARNING] "+format, a...)
}

func (l Logger) Fatal(a ...any) {
	fmt.Println("[FATAL]", a)
}

func (l Logger) Fatalf(format string, a ...any) {
	fmt.Printf("[FATAL] "+format, a...)
}

// 提供一个全局实例，可以直接使用 Log.Debug()
var Log = Logger{}
