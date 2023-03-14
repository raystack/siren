package e2e_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/goto/salt/db"
	"github.com/goto/salt/dockertestx"
	"github.com/goto/salt/log"
	"github.com/goto/siren/config"
	"github.com/goto/siren/internal/store/postgres/migrations"
	"github.com/goto/siren/plugins/providers/cortex"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type CortexTest struct {
	PGConfig          db.Config
	CortexConfig      cortex.AppConfig
	NginxHost         string
	bridgeNetworkName string
	pool              *dockertest.Pool
	network           *dockertest.Network
	resources         []*dockertest.Resource
	hostIP            string
}

func bootstrapCortexTestData(s *suite.Suite, ctx context.Context, client sirenv1beta1.SirenServiceClient, cortexProviderHost string) {
	// add provider cortex
	_, err := client.CreateProvider(ctx, &sirenv1beta1.CreateProviderRequest{
		Host: fmt.Sprintf("http://%s", cortexProviderHost),
		Urn:  "cortex-test",
		Name: "cortex-test",
		Type: "cortex",
	})
	s.Require().NoError(err)

	// add namespace gotocompany-test
	_, err = client.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
		Name:     "fake",
		Urn:      "fake",
		Provider: 1,
	})
	s.Require().NoError(err)

	// add receiver gotocompany-http
	configs, err := structpb.NewStruct(map[string]interface{}{
		"url": "http://fake-webhook-endpoint.gotocompany.com",
	})
	s.Require().NoError(err)
	_, err = client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
		Name: "gotocompany-http",
		Type: "http",
		Labels: map[string]string{
			"entity": "gotocompany",
			"kind":   "http",
		},
		Configurations: configs,
	})
	s.Require().NoError(err)

	// validate
	pRes, err := client.ListProviders(ctx, &sirenv1beta1.ListProvidersRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(pRes.GetProviders()))

	nRes, err := client.ListNamespaces(ctx, &sirenv1beta1.ListNamespacesRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(nRes.GetNamespaces()))

	rRes, err := client.ListReceivers(ctx, &sirenv1beta1.ListReceiversRequest{})
	s.Require().NoError(err)
	s.Require().Equal(1, len(rRes.GetReceivers()))
}

func fetchCortexRules(cortexHost, tenant string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/api/v1/rules", cortexHost), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Scope-OrgID", tenant)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func fetchCortexAlertmanagerConfig(cortexAMHost, tenant string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/api/v1/alerts", cortexAMHost), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Scope-OrgID", tenant)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func triggerCortexAlert(cortexAMHost, tenant, bodyJson string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/api/prom/alertmanager/api/v1/alerts", cortexAMHost), bytes.NewBufferString(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Scope-OrgID", tenant)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func InitCortexEnvironment(appConfig *config.Config) (*CortexTest, error) {
	var (
		err    error
		logger = log.NewZap()
	)

	ct := &CortexTest{
		bridgeNetworkName: fmt.Sprintf("bridge-%s", uuid.New().String()),
		resources:         []*dockertest.Resource{},
	}

	ct.pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	// Create a bridge network for testing.
	ct.network, err = ct.pool.CreateNetwork(ct.bridgeNetworkName)
	if err != nil {
		return nil, err
	}

	dockerNetworks, err := ct.pool.Client.ListNetworks()
	if err != nil {
		return nil, err
	}

	if len(dockerNetworks) == 0 {
		return nil, errors.New("no docker network found")
	}

	for _, ntwrk := range dockerNetworks {
		if ntwrk.Name == "bridge" {
			if len(ntwrk.IPAM.Config) == 0 {
				return nil, errors.New("no docker network gateway config found")
			}
			ct.hostIP = ntwrk.IPAM.Config[0].Gateway
			break
		}
	}

	// pg 1
	logger.Info("creating main postgres...")
	dockerPG, err := dockertestx.CreatePostgres(
		dockertestx.PostgresWithLogger(logger),
		dockertestx.PostgresWithDockertestNetwork(ct.network),
		dockertestx.PostgresWithDockerPool(ct.pool),
		dockertestx.PostgresWithVersionTag("13"),
	)
	if err != nil {
		return nil, err
	}
	ct.resources = append(ct.resources, dockerPG.GetResource())
	logger.Info("main postgres is created")

	logger.Info("creating minio...")
	dockerMinio, err := dockertestx.CreateMinio(
		dockertestx.MinioWithDockertestNetwork(ct.network),
		dockertestx.MinioWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}
	ct.resources = append(ct.resources, dockerMinio.GetResource())
	logger.Info("minio is created")

	logger.Info("migrating minio...")
	if err := dockertestx.MigrateMinio(dockerMinio.GetInternalHost(), "cortex",
		dockertestx.MigrateMinioWithDockertestNetwork(ct.network),
		dockertestx.MigrateMinioWithDockerPool(ct.pool),
	); err != nil {
		return nil, err
	}
	logger.Info("minio is migrated")
	minioURL := fmt.Sprintf("http://%s", dockerMinio.GetInternalHost())

	logger.Info("starting up cortex-am...")
	dockerCortexAM, err := dockertestx.CreateCortex(
		dockertestx.CortexWithModule("alertmanager"),
		dockertestx.CortexWithS3Endpoint(minioURL),
		dockertestx.CortexWithDockertestNetwork(ct.network),
		dockertestx.CortexWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}

	ct.resources = append(ct.resources, dockerCortexAM.GetResource())
	logger.Info("cortex-am is up")

	logger.Info("starting up cortex-all...")
	alertManagerURL := fmt.Sprintf("http://%s/api/prom/alertmanager/", dockerCortexAM.GetInternalHost())
	dockerCortexAll, err := dockertestx.CreateCortex(
		dockertestx.CortexWithAlertmanagerURL(alertManagerURL),
		dockertestx.CortexWithS3Endpoint(minioURL),
		dockertestx.CortexWithDockertestNetwork(ct.network),
		dockertestx.CortexWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}

	ct.resources = append(ct.resources, dockerCortexAll.GetResource())
	logger.Info("cortex-all is up")

	dockerNginx, err := dockertestx.CreateNginx(
		dockertestx.NginxWithDockertestNetwork(ct.network),
		dockertestx.NginxWithDockerPool(ct.pool),
		dockertestx.NginxWithPresetConfig("cortex"),
		dockertestx.NginxWithExposedPort("9009"),
		dockertestx.NginxWithConfigVariables(map[string]string{
			"ExposedPort":      "9009",
			"RulerHost":        dockerCortexAll.GetInternalHost(),
			"AlertManagerHost": dockerCortexAM.GetInternalHost(),
		}),
	)
	if err != nil {
		return nil, err
	}

	ct.NginxHost = dockerNginx.GetExternalHost()
	ct.resources = append(ct.resources, dockerNginx.GetResource())

	ct.PGConfig = db.Config{
		Driver:          "postgres",
		URL:             dockerPG.GetExternalConnString(),
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxLifeTime: time.Millisecond * 100,
		MaxQueryTimeout: time.Millisecond * 100,
	}

	appConfig.DB = ct.PGConfig

	logger.Info("migrating siren...")
	if err = db.RunMigrations(db.Config{
		Driver: appConfig.DB.Driver,
		URL:    appConfig.DB.URL,
	}, migrations.FS, migrations.ResourcePath); err != nil {
		return nil, err
	}
	logger.Info("siren db is migrated")

	return ct, nil
}

func (ct *CortexTest) CleanUp() error {
	for _, r := range ct.resources {
		if err := r.Close(); err != nil {
			return fmt.Errorf("could not purge resource: %w", err)
		}
	}
	if err := ct.network.Close(); err != nil {
		return err
	}
	return nil
}

// func RunDockerizedSiren(apiPort int, pool *dockertest.Pool, netwrk *dockertest.Network) (*dockertest.Resource, error) {
// 	// Build and run the given Dockerfile
// 	resource, err := pool.BuildAndRunWithOptions("../../Dockerfile.dev", &dockertest.RunOptions{
// 		Name:         "siren-e2e",
// 		Networks:     []*dockertest.Network{netwrk},
// 		ExposedPorts: []string{"8080"},
// 		PortBindings: map[docker.Port][]docker.PortBinding{
// 			"8080": []docker.PortBinding{{HostPort: fmt.Sprintf("%d", apiPort)}},
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = pool.Retry(func() error {
// 		var err error
// 		resp, err := http.DefaultClient.Get(fmt.Sprintf("http://localhost:%d", apiPort))
// 		if err != nil {
// 			return err
// 		}
// 		_, err = io.ReadAll(resp.Body)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}

// 	return resource, nil
// }
