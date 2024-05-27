package ceph

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

const (
	EnvDemoUID       = "CEPH_DEMO_UID"
	EnvDemoBucket    = "CEPH_DEMO_BUCKET"
	EnvDemoAccessKey = "CEPH_DEMO_ACCESS_KEY"
	EnvDemoSecretKey = "CEPH_DEMO_SECRET_KEY"
	EnvPublicNetwork = "CEPH_PUBLIC_NETWORK"
	EnvMonitorIP     = "MON_IP"
	EnvRgwName       = "RGW_NAME"
)

const (
	defaultImage = "quay.io/ceph/demo:latest-quincy"

	defaultRgwAccessKey = "demo"
	defaultRgwSecretKey = "b36361c4-1589-42f7-a369-d9dafb926d55"

	defaultRgwPort     = "8080/tcp"
	defaultMonitorPort = "3300/tcp"

	defaultUID    = "demo"
	defaultBucket = "demo"

	defaultPublicNetwork = "0.0.0.0/0"
	defaultMonitorIP     = "127.0.0.1"
	defaultRgwName       = "localhost"

	startRegexFormat = `.*Bucket 's3://%s/' created\n.*`
)

type Container struct {
	testcontainers.Container
	accessKey string
	secretKey string
	bucket    string
}

func (c *Container) GetAccessKey() string {
	return c.accessKey
}

func (c *Container) GetSecretKey() string {
	return c.secretKey
}

func (c *Container) HttpURL(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, defaultRgwPort, "http")
}

func (c *Container) HttpsURL(ctx context.Context) (string, error) {
	return c.PortEndpoint(ctx, defaultRgwPort, "https")
}

func (c *Container) MustHttpURL(ctx context.Context) string {
	url, err := c.HttpURL(ctx)
	if err != nil {
		panic(err)
	}
	return url
}

func (c *Container) MustHttpsURL(ctx context.Context) string {
	url, err := c.HttpsURL(ctx)
	if err != nil {
		panic(err)
	}
	return url
}

func RunContainer(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image: defaultImage,
		Env: map[string]string{
			EnvDemoUID:       defaultUID,
			EnvDemoBucket:    defaultBucket,
			EnvDemoAccessKey: defaultRgwAccessKey,
			EnvDemoSecretKey: defaultRgwSecretKey,
			EnvPublicNetwork: defaultPublicNetwork,
			EnvMonitorIP:     defaultMonitorIP,
			EnvRgwName:       defaultRgwName,
		},
		ExposedPorts: []string{defaultRgwPort, defaultMonitorPort},
	}

	genericContainerReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	}

	for _, opt := range opts {
		if err := opt.Customize(&genericContainerReq); err != nil {
			return nil, err
		}
	}

	if genericContainerReq.WaitingFor == nil {
		regex := fmt.Sprintf(startRegexFormat, genericContainerReq.Env[EnvDemoBucket])
		genericContainerReq.WaitingFor = wait.ForLog(regex).AsRegexp().WithStartupTimeout(5 * time.Minute)
	}

	container, err := testcontainers.GenericContainer(ctx, genericContainerReq)
	if err != nil {
		return nil, err
	}

	accessKey := req.Env[EnvDemoAccessKey]
	secretKey := req.Env[EnvDemoSecretKey]
	bucket := req.Env[EnvDemoBucket]

	result := &Container{Container: container, accessKey: accessKey, secretKey: secretKey, bucket: bucket}
	return result, nil
}
