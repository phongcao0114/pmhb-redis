package klog

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkLogByGlobalLogrus(b *testing.B) {
	b.StopTimer()
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(ioutil.Discard)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		logrus.Infof("lorem ipsum dolor sit amet: %v - %s",
			[]int{0, 1, 2, 3, 4, 5, 6, 7}, "hello world!")
	}
}

func BenchmarkLogByNewKloggerWithMeta(b *testing.B) {
	b.StopTimer()
	bl, err := New(&Config{
		Level:  "debug",
		Format: "text",
		Output: "discard", // file://klog_bench.log
	})
	if err != nil {
		b.Error(err)
	}
	b.StartTimer()

	// TODO: bl == nil after the first iteration if we set output to file, but why?
	// The same code run normally if we use global klog.WithFields function or run it in test
	// func instead of benchmark func (see TestLogByNewKloggerWithMeta below).
	for i := 0; i < b.N; i++ {
		l := bl.
			WithPrefix("bench").
			WithFields(map[string]interface{}{
				"path":       "/api/v1/ping/pong/pem",
				"request_id": "4uenb4syr9u6e544uopr4uenb4syr9u6e544",
			})
		l.Infof("lorem ipsum dolor sit amet: %v - %s",
			[]int{0, 1, 2, 3, 4, 5, 6, 7}, "hello world!")
	}
}

func TestLogByNewKloggerWithMeta(t *testing.T) {
	bl, err := New(&Config{
		Level:  "debug",
		Format: "text",
		Output: "file://klog_test.log",
	})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 1000; i++ {
		l := bl.
			WithPrefix("test").
			WithFields(map[string]interface{}{
				"idx":        i,
				"path":       "/api/v1/ping/pong/pem",
				"request_id": "4uenb4syr9u6e544uopr4uenb4syr9u6e544",
			})
		l.Infof("lorem ipsum dolor sit amet: %v - %s",
			[]int{0, 1, 2, 3, 4, 5, 6, 7}, "hello world!")
	}
}
