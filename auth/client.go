package auth

import (
	"errors"

	gosu "github.com/KommToby/gosu/src"
)

var GosuClient *gosu.GosuClient

func InitClient(clientSecret, clientID string) error {
	if GosuClient != nil {
		return errors.New("gosu client already initialized")
	}

	client, err := gosu.CreateGosuClient(clientSecret, clientID)
	if err != nil {
		return err
	}

	GosuClient = client
	return nil
}
