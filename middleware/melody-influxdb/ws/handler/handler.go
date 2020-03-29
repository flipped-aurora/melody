package handler

import "math/rand"

func PushTestArray() []int {
	var array []int
	for i := 0; i < 7; i++ {
		randNum := rand.Int() % 1500
		array = append(array, randNum)
	}
	return array
}
