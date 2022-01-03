# Bulk Rule management

For org wide use cases, teams end up managing a lot of rules, often manually.

Siren CLI can be used with some Gitops setup to automate the rule creation, rule update, template update. By putting all
the rules and templates YAML files in a version controlled repository, and uploading them using CI Jobs, you get speed
in managing hundreds and thousands of rules in a reliable and predictable manner.

The benefits that one gets via this is:

1. Predictable state of alerts after each CI job run
2. Easy to rollback if something goes wrong
3. Version controlled alerting state, democratizing alert setup, removing dependency from a central team

**Example setup**

1. Create a github repo, let's call it `rules`.
2. Let's create a directory inside it and call it `templates`. This is where people will put the YAML files of
   Templates.
3. Let's create a template names `cpu.yaml` and add the below content
    ```yaml
    apiVersion: v2
    type: template
    name: CPU
    body:
      - alert: CPUWarning
        expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.warning]]
        for: "[[.for]]"
        labels:
          severity: WARNING
        annotations:
          description: CPU has been above [[.warning]] for last [[.for]] {{ $labels.host }}
      - alert: CPUCritical
        expr: avg by (host) (cpu_usage_user{cpu="cpu-total"}) > [[.critical]]
        for: "[[.for]]"
        labels:
          severity: CRITICAL
        annotations:
          description: CPU has been above [[.critical]] for last [[.for]] {{ $labels.host }}
    variables:
      - name: for
        type: string
        default: 10m
        description: For eg 5m, 2h; Golang duration format
      - name: warning
        type: int
        default: 80
      - name: critical
        type: int
        default: 90
    tags:
      - systems
    ```
4. Let's define a shell script which iterates over all files inside `templates/` directory on github to upload templates
   to Siren.

   ```shell
   #!/bin/bash
   set -e
   echo "------------------------------------------------------------"
   echo "Uploading templates: $DIR"
   echo "------------------------------------------------------------"
   
   for FILE in templates/*; do
     eval ./siren template upload $FILE
     echo $'\n'
   done
   
   ```

5. Now as the last step we need to run this script using github action. Here we are pulling siren image and using
   the `upload` command to upload the templates to Siren Web service, denoted by `SIREN_SERVICE_HOST` environment
   variable. An example is:

```yaml
// to be filled later
```

6. For rules, create a directory called `rules` beside `templates` and start define an example rule as given below.
    ```yaml
    apiVersion: v2
    type: rule
    namespace: demo
    provider: production-cortex
    providerNamespace: odpf
    rules:
      TestGroup:
        template: CPU
        status: enabled
        variables:
          - name: for
            value: 15m
          - name: warning
            value: 185
          - name: critical
            value: 195
    ```

7. We can upload the files inside `rules` directory iteratively. Here is an example script. This can be called in github
   CI action.

   ```shell
   #!/bin/bash
   set -e
   echo "------------------------------------------------------------"
   echo "Uploading rules: $DIR"
   echo "------------------------------------------------------------"
   
   for FILE in rules/*; do
     eval ./siren rule upload $FILE
     echo $'\n'
   done
   ```