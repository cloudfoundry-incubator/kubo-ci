package vsphere

import (
	"context"
	"net/url"

	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
)

//go:generate counterfeiter ./ Client
type Client interface {
	DeleteVM(string) error
}

//go:generate counterfeiter ./ vmFinder
type vmFinder interface {
	FindByIp(context.Context, *object.Datacenter, string, bool) (object.Reference, error)
}

//go:generate counterfeiter ./ VM
type VM interface {
	PowerOff() error
	Destroy() error
}

type client struct {
	finder    vmFinder
	extractor func(object.Reference) (VM, error)
}

type internalVM struct {
	vm *object.VirtualMachine
}

func NewClient(vsphereURL *url.URL) (Client, error) {
	ctx := context.Background()
	finder, err := buildSearchIndex(ctx, vsphereURL)
	if err != nil {
		return nil, err
	}
	return NewClientWithFinder(finder, extractVMReference), nil
}

func NewClientWithFinder(finder vmFinder, extractor func(object.Reference) (VM, error)) Client {
	return &client{finder: finder, extractor: extractor}
}

func extractVMReference(r object.Reference) (VM, error) {
	vm, converted := r.(*object.VirtualMachine)
	if !converted {
		return nil, fmt.Errorf("The returned object is not a VM %#v", r)
	}
	return internalVM{vm: vm}, nil
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
	if err != nil || vmReference == nil {
		fmt.Println("VM with IP " + ip + " does not exist")

		return nil
	}
	vm, err := c.extractor(vmReference)
	if err != nil {
		return err
	}

	fmt.Println("Deleting VM with IP " + ip)
	err = vm.PowerOff()
	if err != nil {
		return err
	}
	return vm.Destroy()
}

func (v internalVM) PowerOff() error {
	ctx := context.Background()
	state, err := v.vm.PowerOff(ctx)
	if err != nil {
		return err
	}

	return state.Wait(ctx)
}
func (v internalVM) Destroy() error {
	ctx := context.Background()
	state, err := v.vm.Destroy(ctx)
	if err != nil {
		return err
	}

	return state.Wait(ctx)
}
