package sms

import "errors"

var (
	providers   = map[string]Sender{}
	ErrNotFound = errors.New(`not found provider`)
)

func Clear() {
	providers = map[string]Sender{}
}

func Register(providerName string, sender Sender) {
	providers[providerName] = sender
}

func Provider(providerName string) (Sender, error) {
	if sender, ok := providers[providerName]; ok {
		return sender, nil
	}
	return nil, ErrNotFound
}

func AnyOne() (Sender, string) {
	for providerName, sender := range providers {
		return sender, providerName
	}
	return nil, ``
}
