package address

import "context"

type Info struct {
	City string `json:"City"`
}

type AddressFetcher interface {
	GetByZipCode(ctx context.Context, zipCode string) (*Info, error)
}
