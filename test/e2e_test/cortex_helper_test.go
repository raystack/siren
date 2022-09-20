package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/odpf/salt/db"
	"github.com/odpf/salt/dockertest"
	"github.com/odpf/salt/log"
	"github.com/odpf/siren/config"
	"github.com/odpf/siren/internal/store/postgres/migrations"
	"github.com/odpf/siren/plugins/providers/cortex"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	orydockertest "github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/structpb"
)

type CortexTest struct {
	PGConfig          db.Config
	CortexConfig      cortex.Config
	CortexAMHost      string
	CortexAllHost     string
	bridgeNetworkName string
	pool              *orydockertest.Pool
	network           *docker.Network
	resources         []*orydockertest.Resource
}

func bootstrapCortexTestData(s *suite.Suite, ctx context.Context, client sirenv1beta1.SirenServiceClient, cortexProviderHost string) {
	// add provider cortex
	_, err := client.CreateProvider(ctx, &sirenv1beta1.CreateProviderRequest{
		Host: fmt.Sprintf("http://%s", cortexProviderHost),
		Urn:  "cortex-all-test",
		Name: "cortex-all-test",
		Type: "cortex",
	})
	s.Require().NoError(err)

	// add namespace odpf-test
	_, err = client.CreateNamespace(ctx, &sirenv1beta1.CreateNamespaceRequest{
		Name:     "fake",
		Urn:      "fake",
		Provider: 1,
	})
	s.Require().NoError(err)

	// add receiver odpf-http
	configs, err := structpb.NewStruct(map[string]interface{}{
		"url": "http://fake-webhook-endpoint.odpf.io",
	})
	s.Require().NoError(err)
	_, err = client.CreateReceiver(ctx, &sirenv1beta1.CreateReceiverRequest{
		Name: "odpf-http",
		Type: "http",
		Labels: map[string]string{
			"entity": "odpf",
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

func fetchCortexRules(cortexHost string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/rules", cortexHost))
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

func fetchCortexAlertmanagerConfig(cortexAMHost string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/alerts", cortexAMHost))
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
		resources:         []*orydockertest.Resource{},
	}

	ct.pool, err = orydockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	// Create a bridge network for testing.
	ct.network, err = ct.pool.Client.CreateNetwork(docker.CreateNetworkOptions{
		Name: ct.bridgeNetworkName,
	})
	if err != nil {
		return nil, err
	}

	// pg 1
	logger.Info("creating main postgres...")
	dockerPG, err := dockertest.CreatePostgres(
		dockertest.PostgresWithLogger(logger),
		dockertest.PostgresWithDockerNetwork(ct.network),
		dockertest.PostgresWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}
	ct.resources = append(ct.resources, dockerPG.GetResource())
	logger.Info("main postgres is created")

	logger.Info("creating minio...")
	dockerMinio, err := dockertest.CreateMinio(
		dockertest.MinioWithDockerNetwork(ct.network),
		dockertest.MinioWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}
	ct.resources = append(ct.resources, dockerMinio.GetResource())
	logger.Info("minio is created")

	logger.Info("migrating minio...")
	if err := dockertest.MigrateMinio(dockerMinio.GetInternalHost(), "cortex",
		dockertest.MigrateMinioWithDockerNetwork(ct.network),
		dockertest.MigrateMinioWithDockerPool(ct.pool),
	); err != nil {
		return nil, err
	}
	logger.Info("minio is migrated")
	minioURL := fmt.Sprintf("http://%s", dockerMinio.GetInternalHost())

	logger.Info("starting up cortex-am...")
	dockerCortexAM, err := dockertest.CreateCortex(
		dockertest.CortexWithModule("alertmanager"),
		dockertest.CortexWithS3Endpoint(minioURL),
		dockertest.CortexWithDockerNetwork(ct.network),
		dockertest.CortexWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}
	ct.CortexAMHost = dockerCortexAM.GetExternalHost()
	ct.resources = append(ct.resources, dockerCortexAM.GetResource())
	logger.Info("cortex-am is up")

	logger.Info("starting up cortex-all...")
	alertManagerURL := fmt.Sprintf("http://%s/api/prom/alertmanager/", dockerCortexAM.GetInternalHost())
	dockerCortexAll, err := dockertest.CreateCortex(
		dockertest.CortexWithAlertmanagerURL(alertManagerURL),
		dockertest.CortexWithS3Endpoint(minioURL),
		dockertest.CortexWithDockerNetwork(ct.network),
		dockertest.CortexWithDockerPool(ct.pool),
	)
	if err != nil {
		return nil, err
	}
	ct.CortexAllHost = dockerCortexAll.GetExternalHost()
	ct.resources = append(ct.resources, dockerCortexAll.GetResource())
	logger.Info("cortex-all is up")

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
	logger.Info("siren is migrated")

	return ct, nil
}

func (ct *CortexTest) CleanUp() error {
	for _, r := range ct.resources {
		if err := r.Close(); err != nil {
			return fmt.Errorf("could not purge resource: %w", err)
		}
	}
	if err := ct.pool.Client.RemoveNetwork(ct.network.ID); err != nil {
		return err
	}
	return nil
}
