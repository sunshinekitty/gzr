package comms

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ImageManager is an interface for building an application's image
type ImageManager interface {
	Build(...string) error
	Push(string) error
}

// DockerManager implements ImageManager in order to manage images for Docker
type DockerManager struct{}

// NewDockerManager returns a pointer to an initialized DockerManager
func NewDockerManager() *DockerManager {
	return &DockerManager{}
}

// Build takes a series of arguments to be sent to Docker and builds an image
func (docker *DockerManager) Build(args ...string) error {
	tag, err := GetDockerTag()
	if err != nil {
		return err
	}
	args = append([]string{"build", "-t", tag}, args...)
	build := exec.Command("docker", args...)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	err = build.Run()
	if err != nil {
		return err
	}
	return nil
}

// Push takes a name and pushes this to Docker
func (docker *DockerManager) Push(name string) error {
	push := exec.Command("docker", "push", name)
	push.Stdout = os.Stdout
	push.Stderr = os.Stderr

	err := push.Run()
	if err != nil {
		return err
	}
	return nil
}

// GetDockerTag combines a configured Docker repository name, the current workdinig directory,
// the current time, and a git hash to create a Docker tag appropriate to gzr
// Output format: `repository/$CWD:YYYYMMDD.SHORT_HASH`
func GetDockerTag() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	gm := NewLocalGitManager(path)
	hash, err := gm.CommitHash()
	if err != nil {
		return "", err
	}
	imageName := strings.Split(path, "/")
	name := imageName[len(imageName)-1]
	return fmt.Sprintf("%s/%s:%s.%s", viper.GetString("repository"), name, time.Now().Format("20060102"), hash), nil
}