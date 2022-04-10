package builder

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"
)

// stop container
func stopDockerContainer(name, containerID string) (err error) {
	cmd := []string{"docker", "stop", containerID}
	_, err = execCommand(name, cmd)
	return
}

// run container in the background
func startDockerContainer(name string, environments map[string]string, image string) (containerID, stdOutErr string, err error) {
	// get current dir to mount as volume of the docker container
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	// prepare docker run command
	cmd := []string{"docker", "run", "-d", "-i", "--rm", "-v",
		fmt.Sprintf("%s:/app", currentDir), "-w", "/app"}
	for key, value := range environments {
		cmd = append(cmd, "-e", fmt.Sprintf("%s=%s", key, value))
	}
	cmd = append(cmd, image)
	cmd = append(cmd, "sh")
	log.Debug().Str("stage", name).Msg(strings.Join(cmd, " "))

	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return
	}

	stdOutErr = string(output)
	containerID = strings.TrimSpace(stdOutErr)

	return
}

// exec command in running container
func execDockerCommand(name, containerID, command string) (exitCode int, err error) {
	cmd := []string{"docker", "exec", "-i", containerID}
	cmd = append(cmd, strings.Split(command, " ")...)
	exitCode, err = execCommand(name, cmd)
	return
}

// execute a command and capture stdOut and stdErr outputs
func execCommand(name string, cmd []string) (exitCode int, err error) {
	log.Debug().Str("stage", name).Msg(strings.Join(cmd, " "))

	command := exec.Command(cmd[0], cmd[1:]...)

	// prepare log capture
	stdoutR, stdoutW := io.Pipe()
	command.Stdout = stdoutW
	defer stdoutW.Close()

	stderrR, stderrW := io.Pipe()
	command.Stderr = stderrW
	defer stderrW.Close()

	// start command
	if err := command.Start(); err != nil {
		return -1, fmt.Errorf("executing command %s: %v", cmd, err)
	}

	wg := sync.WaitGroup{}

	stdoutLogger := log.With().Str("stage", name).Logger()
	go logOutput(stdoutR, &stdoutLogger, &wg)

	stderrLogger := log.With().Str("stage", name).Logger()
	go logOutput(stderrR, &stderrLogger, &wg)

	wg.Wait()

	// wait for command execution
	if err := command.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			return int(waitStatus), nil
		}
		return -1, fmt.Errorf("waiting for command %s: %v", cmd, err)
	}

	return
}
