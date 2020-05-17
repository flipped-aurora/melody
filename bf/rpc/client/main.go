package main

import "fmt"

func main() {
	bloomFilter, err := New("127.0.0.1:9999")
	if err != nil {
		fmt.Println("err: " + err.Error())
		return
	}
	defer bloomFilter.Close()

	bloomFilter.Add([]byte("BEARER eyJhbGciOiJSUzI1NiIsImtpZCI6Im1lbG9keSJ9.eyJOaWNrTmFtZSI6ImthcmwiLCJSb2xlSWQiOiIyIiwiVVVJRCI6IjZmMTU4NWQ1LTYzNmUtMTFlYS1hZDk3LTI4N2ZjZjEzZjZmNSIsImV4cCI6MTU4OTM2OTA4MSwiaWF0IjoxNTg5MzYxODgxLCJpc3MiOiJNZWxvZHkifQ.Sr_yYN_UN-E4KakGvhqw43YgCMH6uwjrPHZ02cuFNE-PhH-ujLSsKWgQdt9VSGZ_D1YiSqIvjWv5hC-zKsM0CLt-EQ0JyKbEkErSRKs_GmiV1S-8eSbNTKSlWeE8UzfdsQoJwUwFPJb8VzPrQTXEqokMPE6GoAeU1y73jBhPJvAQHmcw3_kbjsoDQ252_UnBujlXc5XEqdSxjYKg85nN5z-Rl28k39KMLf4yhpUfSt2qa2ENkz8RliqlucAPKvOWw03zLC-KRfWQ8ekw1R4rAHWNsL3p0UOYlZCJfNReYBhFuB470-GdP-f54tGNzhh6zPqbA7qBrObdIIPVkCaarA"))
}
