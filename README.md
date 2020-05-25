![log owl header](https://github.com/jz222/logowl/blob/master/assets/header.png?raw=true)

<div align="center">
  <p>
    <h3>Log Owl</h3>
  </p>
  <p>
    <i>Monitor your services and track errors in production. 🚀📈</i>
  </p>
</div>

---

**Related:**

💻 [**Log Owl Client**](https://github.com/jz222/logowl-client)

📡 [**Log Owl NodeJS Adapter**](https://github.com/jz222/logowl-adapter-nodejs)

🌐 [**Log Owl Browser Adapter**](https://github.com/jz222/logowl-adapter-browser)

<br />

### Features

🔥 **Flexible**

- Group events by services
- Easy-to-use adapter
- Custom adapters for any platform and language
- Customizable
- Can be self-hosted

🔥 **Extensive Event Details**

- Platform information
- Detailed evolution
- Stacktrace
- Code Snippets
- Logs
- Metrics
- Individual badges
- Adapter information

🔥 **Aggregated events**

- Live updates
- Same events are aggregated
- Event count
- Evolution preview

🔥 **User management**

- Invite and remove users from your organization

🔥 **Highly scalable infrastructure**

- Containerized backend
- Simple to deploy and scale

## Run Locally

Running the Log Owl Docker image is the easiest way to get started. All you need is Docker and a MongoDB instance. You can install MongoDB locally, run it as a separate container or use a remote instance provided by services like mLab or Atlas.

![docker logowl](https://github.com/jz222/logowl/blob/master/assets/docker.gif?raw=true)

```
docker run \
--env PORT=2800 \
--env SECRET=secret \
--env MONGO_URI=mongodb://admin:password0@ds263108.mlab.com:63108/logowl-test\?retryWrites=false \
--env MONGO_DB_NAME=logowl-test \
--env MAILGUN_PRIVATE_KEY=aaa-aaa-aaa \
--env MAILGUN_DOMAIN=example.com \
--env CLIENT_URL=http://localhost:3000 \
--env MONTHLY_REQUEST_LIMIT=50000 \
--env IS_SELFHOSTED=true \
-p 2800:2800 \
-it \
jz222/logowl:2.0.0
```

| Environment Variable | Description                                                                                                                                                                              |
|----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| PORT                   | Determines the port that the server should listen on.                                                                                                                                    |
| SECRET                 | Secret key that is used to sign JWT's. Make sure to provide a strong key.                                                                                                                |
| MONGO_URI              | The connection string of the MongoDB. Please refer to the [MongoDB documentation](https://docs.mongodb.com/manual/reference/connection-string/) for the format of the connection string. |
| MONGO_DB_NAME          | The name of the actual database.                                                                                                                                                         |
| MAILGUN_PRIVATE_KEY    | Private key for Mailgun. The private key can be found in the Mailgun settings. This environment variable is optional.                                                                    |
| MAILGUN_DOMAIN         | The domain that is connected to your Mailgun account. This environment variable is optional.                                                                                             |
| CLIENT_URL             | The URL of the Log Owl client.                                                                                                                                                         |
| MONTHLY_REQUEST_LIMIT  | Defines the maximum amount of requests tracked per month. If the limit was reached, incoming requests will no longer be tracked.     |
| IS_SELFHOSTED          | Can either be `true` or `false`. If this environment variable is set to `true`, only one organization can be set up.                                                                 |

## Development Setup

Clone the repository and install dependencies with `go get`. After adding an `.env` file that corresponds to the `.example.env` file, you can start the server with `go run cmd/logowl/main.go`.

## Build

To build a Docker image run the script `build.dev.sh`. It will create a local Docker image called `logowl` that can be run with the Docker command shown above.

## Register an Error

Use the NodeJS adapter to register errors or build your own. To register an error, send a `POST` request to `/logging/error` with a JSON body like shown below.

```json
{
	"ticket": "2ATNP1AD70",
	"message": "test is not a function",
	"path": "/User/example/server/src/server/server.js",
	"line": "15",
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
		"name": "logowl-adapter-nodejs",
		"type": "nodejs",
		"version": "v0.1.0"
	}
}
```

Please notice that timestamps have to be UTC timestamps in seconds.
