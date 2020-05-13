package main

import "fmt"

func main() {
	bloomFilter, err := New("127.0.0.1:9999")
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	defer bloomFilter.Close()

	bloomFilter.Add([]byte("asdasd"))
	check := bloomFilter.Check([]byte("asdasd"))
	fmt.Println(check)
}
