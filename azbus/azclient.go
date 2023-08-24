package azbus

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

const (
	retryBackoff      = 2 * time.Second
	maxRetryDelay     = 20 * time.Second
	maxRetrytAttempts = 3
)

var (
	// NOTE: you don't need to configure these explicitly if you like the defaults.
	// For more information see:
	//  https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus#RetryOptions
	retryOptions = azservicebus.RetryOptions{
		// MaxRetries specifies the maximum number of attempts a failed operation will be retried
		// before producing an error.
		MaxRetries: maxRetrytAttempts,
		// RetryDelay specifies the initial amount of delay to use before retrying an operation.
		// The delay increases exponentially with each retry up to the maximum specified by MaxRetryDelay.
		RetryDelay: retryBackoff,
		// MaxRetryDelay specifies the maximum delay allowed before retrying an operation.
		// Typically the value is greater than or equal to the value specified in RetryDelay.
		MaxRetryDelay: maxRetryDelay,
	}
)

type AZClient struct {
	// ConnectionString contains all the details necessary to connect,
	// authenticate and authorize a client for communicating with azure servicebus.
	ConnectionString string
	client           *azservicebus.Client
}

func NewAZClient(connectionString string) AZClient {
	return AZClient{ConnectionString: connectionString}
}

// azClient - return the client interface
func (c *AZClient) azClient() (*azservicebus.Client, error) {

	if c.client != nil {
		return c.client, nil
	}

	if c.ConnectionString == "" {
		return nil, fmt.Errorf("failed to create client: config must provide either an account name or a connection string")
	}

	client, err := azservicebus.NewClientFromConnectionString(
		c.ConnectionString,
		&azservicebus.ClientOptions{
			RetryOptions: retryOptions,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed creating new client ConnectionString: %w", NewAzbusError(err))
	}
	c.client = client
	return c.client, nil
}
