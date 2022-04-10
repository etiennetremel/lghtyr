package builder

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewBuilder returns a Builder instance
func NewBuilder(filePath string) (builder Builder, err error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("builder.yaml")
	v.AddConfigPath(".")

	if err = v.ReadInConfig(); err != nil {
		err = fmt.Errorf("couldn't find config file located at %s: %v", filePath, err)
		return
	}

	if err = v.Unmarshal(&builder.Config); err != nil {
		err = fmt.Errorf("couldn't read config file located at %s: %v", filePath, err)
		return
	}

	return
}

// RunJob execute a single job given a job name
func (b *Builder) RunJob(name string) (err error) {
	for _, job := range b.Config.Jobs {
		if job.Name == name {
			err = runJob(job)
			if err != nil {
				err = fmt.Errorf("failed running job %s: %v", job.Name, err)
			}
			return
		}
	}
	return fmt.Errorf("no job found with name %s", name)
}

// Run execute all jobs
func (b *Builder) Run() (err error) {
	for _, job := range b.Config.Jobs {
		err = runJob(job)
		if err != nil {
			return fmt.Errorf("failed running job %s: %v", job.Name, err)
		}
	}
	return
}

func runJob(job Job) (err error) {
	for _, stage := range job.Stages {
		err = runStage(stage)
		if err != nil {
			return fmt.Errorf("failed running stage %s: %v", stage.Name, err)
		}
	}

	return
}

func runStage(stage Stage) (err error) {
	containerID, stdOutErr, err := startDockerContainer(stage.Name, stage.Env, stage.Image)
	if err != nil {
		return fmt.Errorf("starting docker container for stage %s: %v: %v", stage.Name, err, stdOutErr)
	}
	defer stopDockerContainer(stage.Name, containerID)

	for _, command := range stage.Commands {
		exitCode, err := execDockerCommand(stage.Name, containerID, command)
		if err != nil {
			return fmt.Errorf("command %s exited with code %d", command, exitCode)
		}
	}
	return
}
