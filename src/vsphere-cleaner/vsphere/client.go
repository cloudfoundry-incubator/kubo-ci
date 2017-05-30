package vsphere

import (
	"context"
	"errors"
	"net/url"
	"vsphere-cleaner/parser"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

type Client interface {
	DeleteVM(string) error
}

type client struct {
	searchIndex object.SearchIndex
}

func BuildUrl(config parser.VMWareConfig) *url.URL {
	parsedUrl := url.URL{
		Scheme: "https",
		Host:   config.IP,
		Path:   "/sdk",
		User:   url.UserPassword(config.User, config.Password),
	}
	return &parsedUrl
}

func NewClient(vsphereURL *url.URL) (Client, error) {
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, vsphereURL, true)
	if err != nil {
		return &client{}, err
	}
	searchIndex := object.NewSearchIndex(c.Client)
	return &client{searchIndex: *searchIndex}, nil
}

func (c *client) DeleteVM(ip string) error {
	ctx := context.Background()
	vmReference, err := c.searchIndex.FindByIp(ctx, nil, ip, true)
	if err != nil {
		return nil
	}

	vm, converted := vmReference.(*object.VirtualMachine)
	if !converted {
		return errors.New("The returned object is not a VM")
	}

	state, err := vm.PowerOff(ctx)
	err = state.Wait(ctx)

	state, err = vm.Destroy(ctx)
	err = state.Wait(ctx)

	return err
}
