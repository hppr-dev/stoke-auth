package benchmark

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

func BenchmarkParallelSingleLogin(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res, err := http.Post(
				"http://localhost:8080/api/login",
				"application/json",
				bytes.NewBuffer([]byte(`{"username" : "tester", "password": "tester"}`)),
			)
			if err != nil {
				log.Println(err)
			}
			if res.StatusCode != 200 {
				log.Println("non 200")
			}

		}
	})
}
