package broadcast

import (
	"math/rand"
	"testing"
)

func TestLogger(t *testing.T) {
}

func BenchmarkLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logger.Infof("log test %v", i)
	}
}

func BenchmarkLoggerParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := rand.Intn(100) + 1
			logger.Infof("log test %v", m)
		}
	})
}
