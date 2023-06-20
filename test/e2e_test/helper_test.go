package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/salt/log"
	"github.com/goto/siren/cli"
	"github.com/goto/siren/config"
	"github.com/goto/siren/core/rule"
	"github.com/goto/siren/core/template"
	sirenv1beta1 "github.com/goto/siren/proto/gotocompany/siren/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

func uploadTemplate(ctx context.Context, cl sirenv1beta1.SirenServiceClient, filePath string) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}

	var yamlObject struct {
		Type string `yaml:"type"`
	}
	err = yaml.Unmarshal(yamlFile, &yamlObject)
	if err != nil {
		return err
	}

	if strings.ToLower(yamlObject.Type) != "template" {
		return errors.New("yaml is not template type")
	}

	var t template.TemplateFile
	err = yaml.Unmarshal(yamlFile, &t)
	if err != nil {
		return err
	}
	body, err := yaml.Marshal(t.Body)
	if err != nil {
		return err
	}

	variables := make([]*sirenv1beta1.TemplateVariables, 0)
	for _, variable := range t.Variables {
		variables = append(variables, &sirenv1beta1.TemplateVariables{
			Name:        variable.Name,
			Type:        variable.Type,
			Default:     variable.Default,
			Description: variable.Description,
		})
	}

	_, err = cl.UpsertTemplate(context.Background(), &sirenv1beta1.UpsertTemplateRequest{
		Name:      t.Name,
		Body:      string(body),
		Variables: variables,
		Tags:      t.Tags,
	})
	if err != nil {
		return err
	}

	return nil
}

func uploadRule(ctx context.Context, cl sirenv1beta1.SirenServiceClient, filePath string) error {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return err
	}

	var yamlObject struct {
		Type string `yaml:"type"`
	}
	err = yaml.Unmarshal(yamlFile, &yamlObject)
	if err != nil {
		return err
	}

	if strings.ToLower(yamlObject.Type) != "rule" {
		return errors.New("yaml is not rule type")
	}

	var yamlBody rule.RuleFile
	err = yaml.Unmarshal(yamlFile, &yamlBody)
	if err != nil {
		return err
	}

	for groupName, v := range yamlBody.Rules {
		var ruleVariables []*sirenv1beta1.Variables
		for i := 0; i < len(v.Variables); i++ {
			v := &sirenv1beta1.Variables{
				Name:  v.Variables[i].Name,
				Value: v.Variables[i].Value,
			}
			ruleVariables = append(ruleVariables, v)
		}

		if yamlBody.Provider == "" {
			return errors.New("provider is required")
		}

		if yamlBody.ProviderNamespace == "" {
			return errors.New("provider namespace is required")
		}

		providersData, err := cl.ListProviders(context.Background(), &sirenv1beta1.ListProvidersRequest{
			Urn: yamlBody.Provider,
		})
		if err != nil {
			return err
		}

		if providersData.GetProviders() == nil {
			return errors.New("provider not found")
		}

		var provider *sirenv1beta1.Provider
		providers := providersData.GetProviders()
		if len(providers) != 0 {
			provider = providers[0]
		} else {
			return errors.New("provider not found")
		}

		res, err := cl.ListNamespaces(context.Background(), &sirenv1beta1.ListNamespacesRequest{})
		if err != nil {
			return err
		}

		if res.GetNamespaces() == nil {
			return errors.New("no response of getting list of namespaces from server")
		}

		var providerNamespace *sirenv1beta1.Namespace
		for _, ns := range res.GetNamespaces() {
			if ns.GetUrn() == yamlBody.ProviderNamespace && ns.Provider == provider.Id {
				providerNamespace = ns
				break
			}
		}

		if providerNamespace == nil {
			return fmt.Errorf("no namespace found with urn: %s under provider %s", yamlBody.ProviderNamespace, provider.Name)
		}

		payload := &sirenv1beta1.UpdateRuleRequest{
			GroupName:         groupName,
			Namespace:         yamlBody.Namespace,
			Template:          v.Template,
			Variables:         ruleVariables,
			ProviderNamespace: providerNamespace.Id,
			Enabled:           v.Enabled,
		}

		_, err = cl.UpdateRule(context.Background(), payload)
		if err != nil {
			fmt.Println(fmt.Sprintf("rule %s/%s/%s upload error",
				payload.Namespace, payload.GroupName, payload.Template), err)
			return err
		}
	}

	return nil
}

func createConnection(ctx context.Context, host string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	return grpc.DialContext(ctx, host, opts...)
}

func CreateClient(ctx context.Context, host string) (sirenv1beta1.SirenServiceClient, func(), error) {
	conn, err := createConnection(context.Background(), host)
	if err != nil {
		return nil, nil, err
	}

	cancel := func() {
		conn.Close()
	}

	client := sirenv1beta1.NewSirenServiceClient(conn)
	return client, cancel, nil
}

func diffYaml(yaml1 []byte, yaml2 []byte) string {
	data1 := make(map[string]any)
	data2 := make(map[string]any)

	err := yaml.Unmarshal(yaml1, &data1)
	if err != nil {
		return "cannot unmarshal yaml1"
	}
	err = yaml.Unmarshal(yaml2, &data2)
	if err != nil {
		return "cannot unmarshal yaml2"
	}

	return cmp.Diff(data1, data2)
}

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func StartSirenServer(cfg config.Config) {
	logger := log.NewZap()
	logger.Info("starting up siren...")
	go func() {
		if err := cli.StartServer(context.Background(), cfg); err != nil {
			logger.Fatal(err.Error())
		}
	}()
	logger.Info("siren is up")

}

func StartSirenMessageWorker(cfg config.Config, closeChannel chan struct{}) error {
	logger := log.NewZap()
	logger.Info("starting up siren notification message worker...")

	if err := cli.StartNotificationHandlerWorker(context.Background(), cfg, closeChannel); err != nil {
		return err
	}
	logger.Info("siren notification message is running")
	return nil
}

func StartSirenDLQWorker(cfg config.Config, closeChannel chan struct{}) error {
	logger := log.NewZap()
	logger.Info("starting up siren notification dlq worker...")

	if err := cli.StartNotificationDLQHandlerWorker(context.Background(), cfg, closeChannel); err != nil {
		return err
	}

	logger.Info("siren notification dlq is running")
	return nil
}
