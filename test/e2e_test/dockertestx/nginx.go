package dockertestx

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	nginxDefaultHealthEndpoint = "/healthz"
	nginxDefaultExposedPort    = "8080"
	nginxDefaultVersionTag     = "1.23"
)

var (
	//go:embed configs/nginx/cortex_nginx.conf
	NginxCortexConfig string
)

type dockerNginxOption func(dc *dockerNginx)

// NginxWithHealthEndpoint is an option to assign health endpoint
func NginxWithHealthEndpoint(healthEndpoint string) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.healthEndpoint = healthEndpoint
	}
}

// NginxWithDockerNetwork is an option to assign docker network
func NginxWithDockerNetwork(network *docker.Network) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.network = network
	}
}

// NginxWithVersionTag is an option to assign version tag
// of a `nginx` image
func NginxWithVersionTag(versionTag string) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.versionTag = versionTag
	}
}

// NginxWithDockerPool is an option to assign docker pool
func NginxWithDockerPool(pool *dockertest.Pool) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.pool = pool
	}
}

// NginxWithDockerPool is an option to assign docker pool
func NginxWithExposedPort(port string) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.exposedPort = port
	}
}

func NginxWithPresetConfig(presetConfig string) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.presetConfig = presetConfig
	}
}

func NginxWithConfigVariables(cv map[string]string) dockerNginxOption {
	return func(dc *dockerNginx) {
		dc.configVariables = cv
	}
}

type dockerNginx struct {
	network            *docker.Network
	pool               *dockertest.Pool
	exposedPort        string
	internalHost       string
	externalHost       string
	presetConfig       string
	versionTag         string
	healthEndpoint     string
	configVariables    map[string]string
	dockertestResource *dockertest.Resource
}

// CreateNginx is a function to create a dockerized nginx
func CreateNginx(opts ...dockerNginxOption) (*dockerNginx, error) {
	var (
		err error
		dc  = &dockerNginx{}
	)

	for _, opt := range opts {
		opt(dc)
	}

	name := fmt.Sprintf("nginx-%s", uuid.New().String())

	if dc.pool == nil {
		dc.pool, err = dockertest.NewPool("")
		if err != nil {
			return nil, fmt.Errorf("could not create dockertest pool: %w", err)
		}
	}

	if dc.versionTag == "" {
		dc.versionTag = nginxDefaultVersionTag
	}

	if dc.exposedPort == "" {
		dc.exposedPort = nginxDefaultExposedPort
	}

	if dc.healthEndpoint == "" {
		dc.healthEndpoint = nginxDefaultHealthEndpoint
	}

	runOpts := &dockertest.RunOptions{
		Name:         name,
		Repository:   "nginx",
		Tag:          dc.versionTag,
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", dc.exposedPort)},
	}

	if dc.network != nil {
		runOpts.NetworkID = dc.network.ID
	}

	var confString string
	switch dc.presetConfig {
	case "cortex":
		confString = NginxCortexConfig
	}

	tmpl := template.New("nginx-config")
	parsedTemplate, err := tmpl.Parse(confString)
	if err != nil {
		return nil, err
	}
	var generatedConf bytes.Buffer
	err = parsedTemplate.Execute(&generatedConf, dc.configVariables)
	if err != nil {
		// it is unlikely that the code returns error here
		return nil, err
	}
	confString = generatedConf.String()

	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var (
		confDestinationFolder = fmt.Sprintf("%s/tmp/dockertest-configs/nginx", pwd)
	)

	foldersPath := []string{confDestinationFolder}
	for _, fp := range foldersPath {
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			if err := os.MkdirAll(fp, 0777); err != nil {
				return nil, err
			}
		}
	}

	if err := os.WriteFile(path.Join(confDestinationFolder, "nginx.conf"), []byte(confString), fs.ModePerm); err != nil {
		return nil, err
	}

	dc.dockertestResource, err = dc.pool.RunWithOptions(
		runOpts,
		func(config *docker.HostConfig) {
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
			config.Mounts = []docker.HostMount{
				{
					Target: "/etc/nginx/nginx.conf",
					Source: path.Join(confDestinationFolder, "nginx.conf"),
					Type:   "bind",
				},
			}
		},
	)
	if err != nil {
		return nil, err
	}

	externalPort := dc.dockertestResource.GetPort(fmt.Sprintf("%s/tcp", dc.exposedPort))
	dc.internalHost = fmt.Sprintf("%s:%s", name, dc.exposedPort)
	dc.externalHost = fmt.Sprintf("localhost:%s", externalPort)

	if err = dc.dockertestResource.Expire(120); err != nil {
		return nil, err
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	dc.pool.MaxWait = 60 * time.Second

	if err = dc.pool.Retry(func() error {
		httpClient := &http.Client{}
		res, err := httpClient.Get(fmt.Sprintf("http://localhost:%s%s", externalPort, dc.healthEndpoint))
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("nginx server return status %d", res.StatusCode)
		}

		return nil
	}); err != nil {
		err = fmt.Errorf("could not connect to docker: %w", err)
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return dc, nil
}

// GetPool returns docker pool
func (dc *dockerNginx) GetPool() *dockertest.Pool {
	return dc.pool
}

// GetResource returns docker resource
func (dc *dockerNginx) GetResource() *dockertest.Resource {
	return dc.dockertestResource
}

// GetInternalHost returns internal hostname and port
// e.g. internal-xxxxxx:8080
func (dc *dockerNginx) GetInternalHost() string {
	return dc.internalHost
}

// GetExternalHost returns localhost and port
// e.g. localhost:51113
func (dc *dockerNginx) GetExternalHost() string {
	return dc.externalHost
}
