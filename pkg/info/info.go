/*
Package info implements a library to expose information about our application.
It will be used by our router package to expose an '/info' endpoint
*/
package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/dstroot/utility"
	"github.com/pkg/errors"
)

const (
	DEFAULT_PORT = "8000"
)

var (
	start = time.Now().UTC()

	// BuildTime is a time label of the moment when the binary was built
	BuildTime = "unset"
	// Commit is a last commit hash at the moment when the binary was built
	Commit = "unset"
	// Version is a semantic version of current build
	Version = "unset"
	// Report exposes our metrics
	Report Metrics
)

// Metrics holds our metrics
type Metrics struct {
	HostName  string
	IPAddress string
	Port      string
	Program   string
	BuildTime string
	Commit    string
	Version   string
	GoVersion string
	PID       int
	RunTime   string
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	return port
}

// Init initializes our metrics.
func Init() (err error) {

	Report.PID = os.Getpid()

	// get hostname
	Report.HostName, err = os.Hostname()
	if err != nil {
		return errors.Wrap(err, "hostname unavailable")
	}

	// get IP address
	Report.IPAddress, err = utility.GetLocalIP()
	if err != nil {
		return errors.Wrap(err, "IP unavailable")
	}

	// get port
	Report.Port = getPort()

	path := strings.Split(os.Args[0], "/")
	Report.Program = strings.Title(path[len(path)-1])
	Report.BuildTime = BuildTime
	Report.Commit = Commit
	Report.Version = Version
	Report.GoVersion = runtime.Version()

	return nil
}

// Handler writes a JSON object with the current metrics
func Handler(w http.ResponseWriter, r *http.Request) {
	Report.RunTime = fmt.Sprintf("%v", utility.RoundDuration(time.Since(start), time.Second))

	j, err := json.MarshalIndent(Report, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

// HandlerFunc returns the info HTTP Handler.
func HandlerFunc() http.Handler {
	return http.HandlerFunc(Handler)
}
