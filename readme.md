## Simple Go Webserver

This repo is meant to be a small example of a good starting point for a Go project.  

* Handles templates (Bootstrap 4 setup as an example)
* Handles a graceful shutdown
* Uses Dep for dependency management
* Sets appropriate timeouts on the http server for production use 
* Uses Httprouter for routing and Negroni for middleware

The repo is structured as follows:

```
├── .editorconfig       # Good starting editor settings
├── .gitattributes      # Good start
├── .gitignore          # Good start
├── Dockerfile          # By placing the Dockerfile on the root
│                       # we can make it very simple - when executing
│                       # the `docker build` command from the root
│                       # the context of the build will carry all the
│                       # relevant Go files and dependencies
├── Gopkg.lock          # Using dep as our package manager
├── Gopkg.toml          # Using dep as our package manager
├── Makefile            # Filled with magic...       
├── pkg                 # Go app packages
│   ├── handlers
|   |   ├── handlers.go
|   |   └── handlers_test.go
│   ├── router
|   |   ├── router.go
|   |   └── router_test.go
│   └── tmpl
|      ├── tmpl.go
|      └── tmpl_test.go
├── VERSION             # What version are we?
├── vendor              # Our vendor libraries
├── main.go             # a `main` file so that one can reference 
└── main_test.go        # the project from the repository name
```

#### References

* [Things to know about HTTP in Go](https://scene-si.org/2017/09/27/things-to-know-about-http-in-go/)
* [What's Coming in Go 1.8](https://tylerchr.blog/golang-18-whats-coming/)
* [So you want to expose Go on the Internet](https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/)
* [Running Go programs on Kubernetes](https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/)
