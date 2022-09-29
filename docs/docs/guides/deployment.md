# Deployment

Siren docker image can be found on Docker hub [here](https://hub.docker.com/r/odpf/siren). You can run the image with
its dependencies.

The dependencies are:

1. Postgres DB
2. Cortex Ruler
3. Cortex Alertmanager

Make sure you have the instances running for them.

## Deploying to Kubernetes

- Create a siren deployment using the helm chart available [here](https://github.com/odpf/charts/tree/main/stable/siren)

# Monitoring

Siren comes with New relic monitoring in built, which user can enable from inside the `config.yaml`

```yaml
newrelic:
  enabled: true
  appname: siren
  license: ____LICENSE_STRING_OF_40_CHARACTERS_____
```

If the `enabled` is set to true, with correct `license` key, you will be able to see the API metrics on your New relic
dashboard.

# Troubleshooting

TBA
