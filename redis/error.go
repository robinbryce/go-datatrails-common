package redis

import (
	"errors"
	"fmt"
)

var (
	ErrNoTenantID           = errors.New("no tenant ID supplied")
	ErrRedisClose           = errors.New("redis close error")
	ErrRedisConnect         = errors.New("redis connect error")
	ErrRedisCounterDisabled = errors.New("redis counter disabled")
	ErrRedisDial            = errors.New("redis dial error")
	ErrRedisDo              = errors.New("redis do error")
	ErrRedisSend            = errors.New("redis send error")
)

func NoTenantIDError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrNoTenantID, name, err)
}

func CloseError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrRedisClose, name, err)
}

func ConnectError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrRedisConnect, name, err)
}

func DialError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrRedisDial, name, err)
}

func DisabledError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrRedisCounterDisabled, name, err)
}

func DoError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrRedisDo, name, err)
}

func SendError(err error, name string) error {
	return fmt.Errorf("%s %s: %w", ErrRedisSend, name, err)
}
