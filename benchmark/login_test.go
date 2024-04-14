package benchmark

import (
	"bytes"
	"net/http"
	"testing"
)

func BenchmarkParallelSingleLogin(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			http.Post(
				"http://localhost:8080/api/login",
				"application/json",
				bytes.NewBuffer([]byte(`{"username" : "tester", "password": "tester" }`)),
			)
		}
	})
}
