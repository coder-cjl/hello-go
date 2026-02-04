package main

import (
	"errors"
	"fmt"
	// "os"
	// "time"
	// "golang.org/x/text/message"
)

type Test struct{}

// / 无参数，无返回值
func (t Test) Method0() {}

// / 有参数，无返回值
func (t Test) Method1(a int, b string) {}

// 无参数，有返回值
func (t Test) Method2() int {
	return 42
}

// / 有参数，有返回值
func (t Test) Method3(a int, b string) (int, string) {
	return a * 2, b + "!"
}

type User struct {
	Name string
	Age  int
}

func (u User) Greet(greeting string) string {
	return greeting + ", " + u.Name
}

func (u User) IsAdult() bool {
	return u.Age >= 18
}

type User1 struct {
	id   int
	name string
}

func (self *User1) GetID() int {
	return self.id
}

// func test() {
// 	u1 := User1{id: 1001, name: "Bob"}
// 	id := u1.GetID
// 	fmt.Println("User ID:", id)
// }

func GetCircleArea(radius float64) (area float64, err error) {
	if radius < 0 {
		err = errors.New("半径不能为负数")
		return
	}

	area = 3.14159 * radius * radius
	return
}

// func test() {
// 	area, err := GetCircleArea(5)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	} else {
// 		fmt.Printf("Circle area: %.2f\n", area)
// 	}
// }

type PathError struct {
	path       string
	op         string
	createTime string
	message    string
}

// func (e *PathError) Error() string {
// 	p := message.NewPrinter(message.MatchLanguage("zh"))
// 	return p.Sprintf("路径错误 - 路径: %s, 操作: %s, 创建时间: %s, 信息: %s", e.path, e.op, e.createTime, e.message)
// }

// func Open(filename string) error {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return &PathError{
// 			path:       filename,
// 			op:         "read",
// 			message:    err.Error(),
// 			createTime: fmt.Sprintf("%v", time.Now().Add()),
// 		}
// 	}

// 	defer file.Close()
// 	return nil
// }

type writer interface {
	Write(data []byte) (n int, err error)
}

type closer interface {
	Close() error
}

type readWriterCloser interface {
	writer
	closer
}

type Dog struct{}

func (d Dog) Write(data []byte) (n int, err error) {
	fmt.Println("Dog is writing data:", string(data))
	return len(data), nil
}

func (d Dog) Close() error {
	fmt.Println("Dog is closing.")
	return nil
}

type Cat struct{}

func (c Cat) Write(data []byte) (n int, err error) {
	fmt.Println("Cat is writing data:", string(data))
	return len(data), nil
}

func (c Cat) Close() error {
	fmt.Println("Cat is closing.")
	return nil
}

func useReadWriterCloser(rwc readWriterCloser) {
	data := []byte("Hello, Go!")
	n, err := rwc.Write(data)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	fmt.Printf("Wrote %d bytes\n", n)

	err = rwc.Close()
	if err != nil {
		fmt.Println("Close error:", err)
		return
	}
	fmt.Println("Closed successfully")
}
