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
