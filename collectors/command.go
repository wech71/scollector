package collectors

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/StackExchange/slog"
)

// command executes the named program with the given arguments.
// If it does not exit within 10s, it is terminated.
func command(name string, arg ...string) ([]byte, error) {
	c := exec.Command(name, arg...)
	var b bytes.Buffer
	c.Stdout = &b
	done := make(chan error, 1)
	go func() {
		done <- c.Run()
	}()
	const commandDuration = time.Second * 10
	select {
	case err := <-done:
		return b.Bytes(), err
	case <-time.After(commandDuration):
		// todo: figure out if this can leave the done chan hanging open
		c.Process.Kill()
		return nil, fmt.Errorf("%v killed after %v", name, commandDuration)
	}
}

func readCommand(line func(string), name string, arg ...string) error {
	b, err := command(name, arg...)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		line(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		slog.Infof("%v: %v\n", name, err)
	}
	return nil
}
