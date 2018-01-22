## Simple Go Webserver

This repo is meant to be a small example of a good starting point for a Go project.  

* Handles templates (Bootstrap 4 setup as an example)
* Handles a graceful HTTP shutdown
* Uses [Dep](https://github.com/golang/dep) for dependency management
* Sets appropriate timeouts on the http server for production use 
* Uses [httprouter](https://github.com/julienschmidt/httprouter) for routing 
* Uses [Negroni](https://github.com/urfave/negroni) for middleware
* Has both expvar and pprof integrated for advanced debugging
* Has prometheus metrics integrated
* Has Jaeger tracing integrated
* Has both "healthz" and "readyz" endpoints for kubernetes
* Has an "info" endpoint to provide program information

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


### Starting the Jaeger Collector

To keep this demo simple, we’ll start the all-in-one Jaeger Docker image:

```sh
$ docker run -d \
    -p 5775:5775/udp \
    -p 6831:6831/udp \
    -p 6832:6832/udp \
    -p 5778:5778 \
    -p 16686:16686 \
    -p 14268:14268 \
    --name=jaeger \
    jaegertracing/all-in-one:latest
```

After a few seconds, all Jaeger components should be up and running.
