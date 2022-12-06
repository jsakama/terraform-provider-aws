// Code generated by internal/generate/servicepackagedata/main.go; DO NOT EDIT.

package docdb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/experimental/intf"
)

var spd = &servicePackageData{}

func registerFrameworkDataSourceFactory(factory func(context.Context) (datasource.DataSourceWithConfigure, error)) {
	spd.frameworkDataSourceFactories = append(spd.frameworkDataSourceFactories, factory)
}

func registerFrameworkResourceFactory(factory func(context.Context) (resource.ResourceWithConfigure, error)) {
	spd.frameworkResourceFactories = append(spd.frameworkResourceFactories, factory)
}

func registerSDKDataSourceFactory(typeName string, factory func() *schema.Resource) {
	spd.sdkDataSourceFactories = append(spd.sdkDataSourceFactories, struct {
		TypeName string
		Factory  func() *schema.Resource
	}{TypeName: typeName, Factory: factory})
}

func registerSDKResourceFactory(typeName string, factory func() *schema.Resource) {
	spd.sdkResourceFactories = append(spd.sdkResourceFactories, struct {
		TypeName string
		Factory  func() *schema.Resource
	}{TypeName: typeName, Factory: factory})
}

type servicePackageData struct {
	frameworkDataSourceFactories []func(context.Context) (datasource.DataSourceWithConfigure, error)
	frameworkResourceFactories   []func(context.Context) (resource.ResourceWithConfigure, error)
	sdkDataSourceFactories       []struct {
		TypeName string
		Factory  func() *schema.Resource
	}
	sdkResourceFactories []struct {
		TypeName string
		Factory  func() *schema.Resource
	}
}

func (d *servicePackageData) Configure(ctx context.Context, meta any) error {
	return nil
}

func (d *servicePackageData) FrameworkDataSources(ctx context.Context) []func(context.Context) (datasource.DataSourceWithConfigure, error) {
	return d.frameworkDataSourceFactories
}

func (d *servicePackageData) FrameworkResources(ctx context.Context) []func(context.Context) (resource.ResourceWithConfigure, error) {
	return d.frameworkResourceFactories
}

func (d *servicePackageData) SDKDataSources(ctx context.Context) []struct {
	TypeName string
	Factory  func() *schema.Resource
} {
	return d.sdkDataSourceFactories
}

func (d *servicePackageData) SDKResources(ctx context.Context) []struct {
	TypeName string
	Factory  func() *schema.Resource
} {
	return d.sdkResourceFactories
}

func (d *servicePackageData) ServicePackageName() string {
	return "docdb"
}

var ServicePackageData intf.ServicePackageData = spd
