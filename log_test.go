package log_test

import (
	"bytes"
	"testing"

	"github.com/trevatk/log"
)

func TestDebug(t *testing.T) {

	var buf bytes.Buffer

	opts := []log.LoggerOption{
		log.WithLevel("DEBUG"),
		log.WithName("unit_test"),
		log.WithWriter(&buf),
	}

	logger := log.New(opts...)

	logger.Debug("unformatted")
	logger.Debugf("%d", 1)

	if buf.Len() == 0 {
		t.Fatal("buffer length is empty")
	}

	t.Logf("%s", buf.String())
}

func TestInfo(t *testing.T) {

	var buf bytes.Buffer

	opts := []log.LoggerOption{
		log.WithLevel("INFO"),
		log.WithName("unit_test"),
		log.WithWriter(&buf),
	}

	logger := log.New(opts...)

	logger.Debug("unformatted")

	if buf.Len() != 0 {
		t.Fatalf("buffer length %d greater than expected %d", buf.Len(), 0)
	}

	logger.Info("unformatted")
	logger.Infof("%d", 1)

	if buf.Len() == 0 {
		t.Fatal("buffer length is empty")
	}

	t.Logf("%s", buf.String())
}

func TestWarn(t *testing.T) {

	var buf bytes.Buffer

	opts := []log.LoggerOption{
		log.WithLevel("WARN"),
		log.WithName("unit_test"),
		log.WithWriter(&buf),
	}

	logger := log.New(opts...)

	logger.Info("unformatted")

	if buf.Len() != 0 {
		t.Fatalf("buffer length %d greater than expected %d", buf.Len(), 0)
	}

	logger.Warn("unformatted")
	logger.Warnf("%d", 1)

	if buf.Len() == 0 {
		t.Fatal("buffer length is empty")
	}

	t.Logf("%s", buf.String())
}

func TestError(t *testing.T) {

	var buf bytes.Buffer

	opts := []log.LoggerOption{
		log.WithLevel("ERROR"),
		log.WithName("unit_test"),
		log.WithWriter(&buf),
	}

	logger := log.New(opts...)

	logger.Info("unformatted")

	if buf.Len() != 0 {
		t.Fatalf("buffer length %d greater than expected %d", buf.Len(), 0)
	}

	logger.Error("unformatted")
	logger.Errorf("%d", 1)

	if buf.Len() == 0 {
		t.Fatal("buffer length is empty")
	}

	t.Logf("%s", buf.String())
}
