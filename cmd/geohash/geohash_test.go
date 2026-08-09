package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/pierrre/assert"
)

func captureStdout(t *testing.T, f func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w
	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()
	ferr := f()
	_ = w.Close()
	os.Stdout = old
	<-done
	return buf.String(), ferr
}

func withStdin(t *testing.T, input string, f func() error) error {
	t.Helper()
	old := os.Stdin
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	_, _ = io.Copy(w, bytes.NewReader([]byte(input)))
	_ = w.Close()
	os.Stdin = r
	ferr := f()
	os.Stdin = old
	return ferr
}

func TestProcessStdinGeohash(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return withStdin(t, "u09tvqx", processStdin)
	})
	assert.NoError(t, err)
	assert.StringHasSuffix(t, out, "\n")
}

func TestProcessStdinLatLon(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return withStdin(t, "48.86,2.35", processStdin)
	})
	assert.NoError(t, err)
	assert.StringHasSuffix(t, out, "\n")
}

func TestProcessStdinMultiple(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return withStdin(t, "u09tvqx 48.86,2.35", processStdin)
	})
	assert.NoError(t, err)
	assert.StringHasSuffix(t, out, "\n")
}

func TestProcessStdinEmpty(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return withStdin(t, "", processStdin)
	})
	assert.NoError(t, err)
	assert.Equal(t, out, "")
}

func TestProcessStdinError(t *testing.T) {
	_, err := captureStdout(t, func() error {
		return withStdin(t, "1,2,3", processStdin)
	})
	assert.Error(t, err)
}

func TestProcessArgsGeohash(t *testing.T) {
	flagPrecision = 0
	flagRound = true
	assert.NoError(t, flag.CommandLine.Parse([]string{"u09tvqx"}))
	out, err := captureStdout(t, processArgs)
	assert.NoError(t, err)
	assert.StringHasSuffix(t, out, "\n")
}

func TestProcessArgsError(t *testing.T) {
	flagPrecision = 0
	flagRound = true
	assert.NoError(t, flag.CommandLine.Parse([]string{"1,2,3"}))
	_, err := captureStdout(t, processArgs)
	assert.Error(t, err)
}

func TestProcessSwitchWithArgs(t *testing.T) {
	flagPrecision = 0
	flagRound = true
	assert.NoError(t, flag.CommandLine.Parse([]string{"u09tvqx"}))
	out, err := captureStdout(t, processSwitch)
	assert.NoError(t, err)
	assert.StringHasSuffix(t, out, "\n")
}

func TestProcessSwitchWithStdin(t *testing.T) {
	assert.NoError(t, flag.CommandLine.Parse([]string{}))
	out, err := captureStdout(t, func() error {
		return withStdin(t, "u09tvqx", processSwitch)
	})
	assert.NoError(t, err)
	assert.StringHasSuffix(t, out, "\n")
}

func TestProcessValueInvalidLat(t *testing.T) {
	_, err := processValue("abc,2.35")
	assert.Error(t, err)
}

func TestProcessValueInvalidLon(t *testing.T) {
	_, err := processValue("48.86,abc")
	assert.Error(t, err)
}

func TestProcessValuePrecision(t *testing.T) {
	flagPrecision = 12
	t.Cleanup(func() {
		flagPrecision = 0
	})
	gh, err := processValue("48.86,2.35")
	assert.NoError(t, err)
	assert.Equal(t, gh, "u09tvqxnnuph")
}

func TestProcessValueDecodeError(t *testing.T) {
	_, err := processValue("é")
	assert.Error(t, err)
}

func TestProcessValueNoRound(t *testing.T) {
	flagRound = false
	t.Cleanup(func() {
		flagRound = true
	})
	result, err := processValue("u09tvqx")
	assert.NoError(t, err)
	assert.StringContains(t, result, ",")
}
