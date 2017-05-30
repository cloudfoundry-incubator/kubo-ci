package vsphere

import (
	"context"
	"errors"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"vsphere-cleaner/parser"
)

type Client interface {
	DeleteVM(string) error
}

//go:generate counterfeiter ./ vmFinder
type vmFinder interface {
	FindByIp(context.Context, *object.Datacenter, string, bool) (object.Reference, error)
}

type client struct {
	finder vmFinder
}

func BuildUrl(config parser.VSphereConfig) *url.URL {
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
	finder, err := buildSearchIndex(ctx, vsphereURL)
	if err != nil {
		return nil, err
	}
	return NewClientWithFinder(finder), nil
}

func NewClientWithFinder(finder vmFinder) Client {
	return &client{finder: finder}
}

func buildSearchIndex(ctx context.Context, vsphereURL *url.URL) (vmFinder, error) {
	c, err := govmomi.NewClient(ctx, vsphereURL, true)
	if err != nil {
		return nil, err
	}
	return object.NewSearchIndex(c.Client), nil
}

func (c *client) DeleteVM(ip string) error {
	ctx := context.Background()
	vmReference, err := c.finder.FindByIp(ctx, nil, ip, true)
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
