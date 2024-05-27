package ceph

import "github.com/testcontainers/testcontainers-go"

func WithAccessKey(key string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env[EnvDemoAccessKey] = key
		return nil
	}
}

func WithSecretKey(key string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env[EnvDemoSecretKey] = key
		return nil
	}
}

func WithBucket(name string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Env[EnvDemoBucket] = name
		return nil
	}
}

func WithSSLDisabled() testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.Entrypoint = []string{
			"bash",
			"-c",
			`sed -i '/^rgw frontends = .*/a rgw verify ssl = false\
rgw crypt require ssl = false' /opt/ceph-container/bin/demo;
/opt/ceph-container/bin/demo;`,
		}
		return nil
	}
}
