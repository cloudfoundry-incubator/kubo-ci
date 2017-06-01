package vsphere

import (
	"context"
	"errors"
	"net/url"

	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

//go:generate counterfeiter ./ Client
type Client interface {
	DeleteVM(string) error
}

type vmFinder interface {
	FindByIp(context.Context, *object.Datacenter, string, bool) (object.Reference, error)
}

type client struct {
	finder vmFinder
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
		fmt.Println("VM with IP " + ip + " does not exist")

		return nil
	}

	vm, converted := vmReference.(*object.VirtualMachine)
	if !converted {
		return errors.New("The returned object (IP = '" + ip + "') is not a VM")
	}

	fmt.Println("Deleting VM with IP " + ip)
	state, err := vm.PowerOff(ctx)
	// TEST ME PLIIIZZ!
	err = state.Wait(ctx)

	state, err = vm.Destroy(ctx)
	err = state.Wait(ctx)

	return err
}
