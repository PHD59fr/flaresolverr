# FlareSolverr Go Client 🚀

This Go package facilitates interactions with FlareSolverr, designed to solve JavaScript challenges presented by websites protected with Cloudflare. It manages standard HTTP requests such as GET, POST, PUT, and DELETE, routing them through FlareSolverr to circumvent protections.

## Usage 🛠

### Initializing the Client and Making Requests 🌐

Below is an example of how to create a client instance with your FlareSolverr base URL and any additional configurations, followed by executing a GET or POST request.

```go
package main

import (
	"github.com/phd59fr/flaresolverr"
	"github.com/sirupsen/logrus"
)

func main() {
	// Set up the logrus logger
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	// Initialize the FlareSolverr client
	clientInput := flaresolverr.ClientInput{
		BaseUrl: "http://localhost:8191",
		TimeOut: 60000, // Timeout in milliseconds
	}

	client, err := flaresolverr.Init(clientInput)
	if err != nil {
		logger.Fatalf("Error initializing FlareSolverr: %v", err)
	}

	logger.Infof("FlareSolverr Client successfully initialized. Version: %s", client.Version)

	// Perform a GET request
	getResponse, err := client.Get("https://www.google.com")
	if err != nil {
		logger.Fatalf("Error performing GET request: %v", err)
	}

	logger.Infof("GET request response: %s", getResponse.Solution.Response)

	// Perform a POST request with data
	postData := map[string]interface{}{
		"key": "value",
	}

	postResponse, err := client.Post("https://example.com/post", postData)
	if err != nil {
		logger.Fatalf("Error performing POST request: %v", err)
	}

	logger.Infof("POST request response: %s", postResponse.Solution.Response)
}
```

## Error Handling ⚠️

All functions return an `error` object which should be checked to ensure the request was processed successfully.
