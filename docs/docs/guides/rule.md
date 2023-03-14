import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Rule

export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

Siren rules are generated from predefined [templates](template.md) by providing values of the variables of the template.

One can create templates using either HTTP APIs or CLI.

## API interface

A rule is uniquely identified with the combination of provider's namespace (uniquely identifies which provider and namespace), template name, optional namespace, optional group name.

One can choose any namespace and group name. In Cortex terminology, namespace is a collection of groups. Groups can have
one or more rules.

### Creating/Updating a rule

The below snippet describes an example of rule creation/update. Same API can be used to enable or disable alerts.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren rule create --file rule.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request PUT
  --url `}{defaultHost}{`/`}{apiVersion}{`/rules
  --header 'content-type: application/json'
  --data-raw '{
  "namespace": "gotocompany",
  "group_name": "CPUHigh",
  "template": "CPU",
  "providerNamespace": "3"
  "variables": [
    {
      "name": "for",
      "value": "15m",
      "type": "string"
    },
     {
      "name": "team",
      "value": "gotocompany",
      "type": "string"
    }
  ],
  "enabled": true,
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

Here we are using CPU template and providing value for few variables("for", "team"). In case some variables value is not
provided default will be picked from the template's definition.

### Terminology of the request body

| Term              | Description                                                    | Example/Default   |
| ----------------- | -------------------------------------------------------------- | ----------------- |
| Namespace         | Corresponds to Cortex namespace in which rule will be created  | kafka             |
| Group Name        | Corresponds to Cortex group name in which rule will be created | CPUHigh           |
| providerNamespace | Corresponds to a tenant in a provider                          | 4                 |
| Template          | what template is used to create the rule                       | CPU               |
| Variables         | Value of variables defined inside the template                 | See example above |
| Enabled           | boolean describing if the rule is enabled or not               | true              |

**Fetching rules**

Rules can be fetched and filtered with multiple parameters. An example of all filters is described below.

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren rule list --namespace foo --provider-namespace 4 --group-name CPUHigh --template CPU
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request GET
  --url `}{defaultHost}{`/`}{apiVersion}{`/rules?namespace=foo&providerNamespace=4&groupName=CPUHigh&template=CPU`}
    </CodeBlock>
  </TabItem>
</Tabs>

## CLI Interface

```text
Work with rules.

rules are used for alerting within a provider.

Usage:
  siren rule [command]

Aliases:
  rule, rules

Available Commands:
  edit        Edit a rule
  list        List rules
  upload      Upload Rules YAML file

Flags:
  -h, --help   help for rule

Use "siren rule [command] --help" for more information about a command.
```

With CLI, you will need a YAML file in the below specified format to create/edit rules.
**Example rule file**

```yaml
apiVersion: v2
type: rule
namespace: demo
provider: localhost-cortex
providerNamespace: test
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

In the above example, we are creating rules
inside `demo` [namespace](https://cortexmetrics.io/docs/api/#get-rule-groups-by-namespace)
under `test` [tenant](https://cortexmetrics.io/docs/architecture/#the-role-of-prometheus) of `localhost-cortex`
provider.

The rules array defines actual rules defined over the templates. Here `TestGroup` is the name of the group which will be
created/updated with the rule defined by `CPU` template. The example shows the value of variables provided in creating
rules(alert).

**Example upload command**

```shell
siren rule upload cpu_rule.yaml
```

The yaml file can be edited and re-uploaded to edit the rule thresholds.

### Terminology

| Term              | Description                                                              | Example/Default   |
| ----------------- | ------------------------------------------------------------------------ | ----------------- |
| API Version       | Which API to use to parse the YAML file                                  | v2                |
| Type              | Describes the type of object represented by YAML file                    | rule              |
| Namespace         | Corresponds to Cortex namespace in which rule will be created            | kafka             |
| Entity            | Corresponds to tenant name in cortex                                     | gotocompany              |
| Rules             | Map of GroupNames describing what template is used in a particular group | See example file  |
| Variables         | Value of variables defined inside the template                           | See example above |
| provider          | URN of monitoring provider to be used                                    | localhost-cortex  |
| providerNamespace | URN of tenant to choose inside the monitoring provider                   | test              |

## Managing Templates via YAML File

Siren gives flexibility to templatize rules for re-usability purpose. Template can be managed via APIs (REST
and GRPC). Apart from that, there is a command line interface as well which parses a YAML file in a specified format (as
described below) and upload to Siren using an HTTP Client of Siren Service. Refer [here](../guides/template.md) for
more details around usage and terminology.

## Managing Rules via YAML File

To manage rules in bulk, Siren gives a way to manage rules using YAML files, which you can manage in a version
controlled repository. Using the `upload` command one can upload a rule YAML file in a specified format (as described
below) and upload to Siren using the GRPC Client(comes inbuilt) of Siren Service. Refer [here](../guides/rule.md) for
more details around usage and terminology.

**Note:** Updating a template also updates the associated rules.

# Bulk Rule management

For org wide use cases, teams end up managing a lot of rules, often manually. Siren CLI can be used to automate the rule creation, rule update, and [template](./template.md) update.

## Use Case: CI

The Siren CLI could further be used in GitOps scenario by putting all the rules and templates YAML files in a version controlled repository and uploading them using CI Jobs. By doing so, you will get speed in managing hundreds and thousands of rules in a reliable and predictable manner.

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
     eval $ siren template upload $FILE
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
   providerNamespace: gotocompany
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
     eval $ siren rule upload $FILE
     echo $'\n'
   done
   ```
