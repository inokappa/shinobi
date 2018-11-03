package main

import (
	"bytes"
	"os"
	"os/exec"
    "fmt"
    "strings"
    "time"
    "testing"
)

func setup() {
    fmt.Println("setup...")
}

func teardown() {
    fmt.Println("teardown.")
}

func TestVersionFlag(t *testing.T) {
	cmd := exec.Command("gom", "run", "shinobi.go", "-version")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout

	_ = cmd.Run()

	if ! strings.Contains(stdout.String(), AppVersion) {
        t.Fatal("Failed Test")
    }
}

func TestStdoutList(t *testing.T) {
	cmd := exec.Command("./tests/scripts/test_stdout_list.sh")
    stdout := new(bytes.Buffer)
    cmd.Stdout = stdout
	output := "dummy-user-pool"

    _ = cmd.Run()

    if ! strings.Contains(stdout.String(), output) {
        t.Fatal("Failed Test")
    }
}

func TestStdoutListHasUser(t *testing.T) {
    cmd := exec.Command("./tests/scripts/test_stdout_list_has_user.sh")
    stdout := new(bytes.Buffer)
    cmd.Stdout = stdout
	output := "testtest1"

    _ = cmd.Run()

    if ! strings.Contains(stdout.String(), output) {
        t.Fatal("Failed Test")
    }
}

func TestConvertDate(t *testing.T) {
    str := "2018-09-28 22:52:24 +0000 UTC"
    layout := "2006-01-02 15:04:05 +0000 UTC"
    tm, _ := time.Parse(layout, str)

    actual := convertDate(tm)
    expected := "2018-09-29 07:52:24"

    if actual != expected {
        t.Fatal("Failed Test")
    }
}

func TestMain(m *testing.M) {
    setup()
    ret := m.Run()
    teardown()
    os.Exit(ret)
}