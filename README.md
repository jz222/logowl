# LOGGY - Monitoring Service

LOGGY allows you to monitor your services and track errors in production.

This repository contains the API of the service. Please find the client and the NodeJS adapter in the respective repo.

- [LOGGY Client](https://github.com/jz222/loggy-client)
- [LOGGY NodeJS Adapter](https://github.com/jz222/loggy-adapter-nodejs)

## Run Locally

Running the LOGGY Docker image is the easiest way to get started. All you need is Docker and a MongoDB instance. You can install MongoDB locally, run it as a separate container or use a remote instance provided by services like mLab or Atlas.

![docker loggy](https://github.com/jz222/loggy/blob/master/assets/docker-loggy.gif?raw=true)

```
docker run \
--env PORT=2800 \
--env SECRET=secret \
--env MONGO_URI=mongodb://admin:password0@ds263108.mlab.com:63108/loggy-test\?retryWrites=false \
--env MONGO_DB_NAME=loggy-test \
--env IS_SELFHOSTED=true \
-p 2800:2800 \
-it \
jz222/loggy:0.1.0
```

#### Port

Determines the port that the server should listen on.

#### Secret

Secret key that is used to sign JWT's. Make sure to provide a strong key.

#### MongoDB URI

The connection string of the MongoDB. Please refer to the [MongoDB documentation](https://docs.mongodb.com/manual/reference/connection-string/) for the format of the connection string.

#### MongoDB Database Name

The name of the actual database.

#### Is Selfhosted

Can either be `true` or `false`. If this environment variable is set to `true`, only one organization can be set up.

## Development Setup

Clone the repository and install dependencies with `go get`. After adding an `.env` file that corresponds to the `.example.env` file, you can start the server with `go run main.go`.

## Build

To build a Docker image run the script `build.dev.sh`. It will create a local Docker image called `loggy` that can be run with the Docker command shown above.

## Register an Error

Use the NodeJS adapter to register errors or build your own. To register an error, send a `POST` request to `/logging/error` with a JSON body like shown below.

```json
{
	"ticket": "2ATNP1AD70",
	"message": "test is not a function",
	"path": "/User/example/server/src/server/server.js",
	"line": 15,
	"stacktrace": "the error stack trace",
	"badges": {
		"cluster": "test"
	},
	"type": "exception",
	"metrics": {
		"platform": "linux"
	},
	"logs": [
		{
			"type": "info",
			"log": "process started",
			"timestamp": 1585689440
		}
	],
	"snippet": {
		"10": "        cluster: 'EU',",
        	"11": "        serviceID: '20010-A'",
        	"12": "    }",
        	"13": "});",
        	"14": "",
        	"15": "test();",
        	"16": "",
        	"17": "// Routes",
        	"18": "const routes = require('../routes');",
        	"19": "",
        	"20": "// Configs"
	},
	"timestamp": 1585689898,
	"adapter": {
		"name": "loggy-adapter-nodejs",
		"type": "nodejs",
		"version": "v0.1.0"
	}
}
```

Please notice that timestamps have to be UTC timestamps in seconds.
