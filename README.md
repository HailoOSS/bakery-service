# My Service

This is a template for h2o Go projects.

Simple steps for creating your project:

  1. Create a folder within the `proto` folder for each endpoint that you need
  2. Add a `.proto` file within each folder, defining your `Request` and
     `Response` for each endpoint
  3. Create a `handler` for each endpoint within the `handler` folder
  4. Open up `main.go` and fill in the details about the project (review the
     service name and tier, fill in the description and your contact details)
  5. Register each of your endpoints (again in `main.go`)
  6. Run `go get github.com/hailcab/your-repo-name` so that Go will compile
     the binary (and put it in your path)
  7. Try it out! (it should be in your path)

Once you've got it working, you're ready to "deploy".


## More details

### Proto

The `proto` folder is going to have one folder per endpoint. Within each folder
there should be a `.proto` file that is named the same as the folder. This
should define the `Request` and `Response` format for the endpoint. This
defined your interface - it tells users of your service what stuff they need
to send the service, and what the service will send back.

So if we had a service with two endpoints, `read` and `create`, we would have
the following structure:

	proto/
		read/
			read.proto
		create/
			create.proto

### Handlers

A service is just a set of endpoints. Each endpoint has a protobuf `Request`
and `Response` format, plus a `handler`. The handler is a function that will
be fed requests to respond to. The handler is called in its own "goroutine"
so anything you do or call should be "thread safe".

By convention, handlers live in a sub-package called `handler` and the
function name matches the endpoint name.


