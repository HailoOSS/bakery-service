# My Service

This is a template for h2o Go projects.

Simple steps for creating your project:

  1. Create a folder within the `proto` folder for each endpoint that you need
  2. Add a `.proto` file within each folder, defining your `Request` and
     `Response` for each endpoint
  3. Compile the protobuf code by running `./common/script/protoc.sh` from the
     project root (note that you'll need to `git submodule init` and `git
     submodule update` in order to pick up the common submodule)
  4. Create a `handler` for each endpoint within the `handler` folder
  5. Open up `main.go` and fill in the details about the project (review the
     service name and tier, fill in the description and your contact details)
  6. Register each of your endpoints (again in `main.go`)
  7. Run `go get github.com/hailocab/bakery-service` so that Go will compile
     the binary (and put it in your path)
  8. Try it out! (it should be in your path)

Once you've got it working, you're ready to "deploy". Deployment consists
of two steps -- a **build** step and then a **provisioning** step.

To build:

  1. Push your code to Git
  2. Type `hubot: ci setup hailocab/bakery-service` in IRC

To provision (actually run your service on some servers in an environment) check https://hailo.jira.com/wiki/display/HTWO/2013/10/24/Deploying+a+service+with+hshell+or+the+dashboard

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


