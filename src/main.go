package main

// func test1() {
// 	a := Add(1, 2)
// 	fmt.Println(a)
// }

// func test2() {
// 	a, s := Test1(1, 2, "cjl")
// 	fmt.Println(a, s)
// }

// func test3() {
// 	t := Test{}
// 	t.Method0()
// 	t.Method1(10, "hello")
// 	r2 := t.Method2()
// 	fmt.Println(r2)
// 	r3a, r3b := t.Method3(5, "go")
// 	fmt.Println(r3a, r3b)
// }

// func test4() {
// 	u := User{Name: "Alice", Age: 20}

// 	greeting := u.Greet("Hello")
// 	fmt.Println(greeting)

// 	isAdult := u.IsAdult()
// 	fmt.Println(isAdult)
// }

func test5() {
	http := TestHttp{}
	http.Start()
}

func test6() {
	o := TestGoRoutine{}
	o.Test()
}

func mysqlTest() {
	// t := MyGormSQL{}
	// t.Test()
	// t := GoRedis{}
	// t := GoEtcd{}
	t := GoKafka{}
	t.Test()
}

func main() {
	mysqlTest()
	// test6()
	// test5()
	// test4()
	// test3()
	// test2()
	// test1()
}
