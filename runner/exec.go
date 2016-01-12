package runner

import (
	"io"
	"sync"

	"github.com/drone/drone-go/drone"
	"github.com/samalba/dockerclient"
)

// State represents the state of an execution.
type State struct {
	sync.Mutex

	Repo      *drone.Repo
	Build     *drone.Build
	BuildLast *drone.Build
	Job       *drone.Job
	System    *drone.System
	Workspace *drone.Workspace

	// Client is an instance of the Docker client
	// used to spawn container tasks.
	Client dockerclient.Client

	Stdout, Stderr io.Writer
}

// Exit writes the exit code. A non-zero value
// indicates the build exited with errors.
func (s *State) Exit(code int) {
	s.Lock()
	defer s.Unlock()

	// only persist non-zero exit
	if code != 0 {
		s.Job.ExitCode = code
		s.Job.Status = drone.StatusFailure
		s.Build.Status = drone.StatusFailure
	}
}

// ExitCode reports the process exit code. A non-zero
// value indicates the build exited with errors.
func (s *State) ExitCode() int {
	s.Lock()
	defer s.Unlock()

	return s.Job.ExitCode
}

// Failed reports whether the execution has failed.
func (s *State) Failed() bool {
	return s.ExitCode() != 0
}
