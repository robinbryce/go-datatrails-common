// Package azblob reads/writes files to Azure
// blob storage in Chunks.
package azblob

import (
	"errors"
	"fmt"

	azStorageBlob "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var (
	ErrUnspecifiedContainer = errors.New("storer: container is unspecified")
)

// so we dont have to import azure blob package anywhere else
type ContainerClient = azStorageBlob.ContainerClient
type ServiceClient = azStorageBlob.ServiceClient
type SharedKeyCredential = azStorageBlob.SharedKeyCredential

// Storer implements usage of Reader/Writer backed by azblob
type Storer struct {
	AccountName   string
	ResourceGroup string
	Subscription  string
	Container     string

	credential      *SharedKeyCredential
	rootURL         string
	containerURL    string
	containerClient *ContainerClient
	serviceClient   *ServiceClient

	log                          Logger
	setReadResponseScannedStatus ReadResponseScannedStatus

	startSpanFromContext startSpanFromContextFunc
}

type StorerOption func(*Storer)

func WithSetScannedStatus(s ReadResponseScannedStatus) StorerOption {
	return func(a *Storer) {
		a.setReadResponseScannedStatus = s
	}
}

func WithStorerSpanFromContext(s startSpanFromContextFunc) StorerOption {
	return func(a *Storer) {
		a.startSpanFromContext = s
	}
}

// New returns new az blob read/write object
func New(
	log Logger,
	accountName string,
	resourceGroup string,
	subscription string,
	container string,
	options ...StorerOption,
) (*Storer, error) {

	var err error
	log.Debugf("New Storer: %s/%s/%s/%s",
		accountName,
		resourceGroup,
		subscription,
		container,
	)

	secret, credential, err := credentials(
		accountName,
		resourceGroup,
		subscription,
	)
	if err != nil {
		return nil, err
	}
	rootURL := secret.URL

	if container == "" {
		log.Infof("Storer: %v", ErrUnspecifiedContainer)
		return nil, ErrUnspecifiedContainer
	}
	azp := Storer{
		AccountName:   accountName,
		ResourceGroup: resourceGroup,
		Subscription:  subscription,
		Container:     container,
		credential:    credential,
		rootURL:       rootURL,
		log:           log,
	}
	for _, option := range options {
		option(&azp)
	}

	azp.containerURL = fmt.Sprintf(
		"%s%s",
		rootURL,
		container,
	)
	azp.serviceClient, err = azStorageBlob.NewServiceClientWithSharedKey(
		rootURL,
		credential,
		nil,
	)
	if err != nil {
		log.Infof("unable to create serviceclient %s: %v", azp.containerURL, err)
		return nil, err
	}
	azp.containerClient, err = azp.serviceClient.NewContainerClient(container)
	if err != nil {
		log.Infof("unable to create containerclient %s: %v", container, err)
		return nil, err
	}

	return &azp, nil
}
