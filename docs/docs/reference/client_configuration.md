# Client Configuration

When using `siren` client CLI, sometimes there are client-specifi flags that are required to be passed e.g. `--host` so you are calling Siren like this.

```bash
siren receiver list --host localhost:8080
```
Siren client CLI could use a client config so you don't need to pass client-required flags e.g. `--host` everytime you run `siren` command. For Siren client CLI, here is the required config to interact to Siren server.

```yaml
host: localhost:8080
```
You could easily generate client config by running this command:
```bash
siren config init
```
This will create (if not exists) a config file `${HOME}/.config/gotocompany/siren.yaml` with default values. You can modify the value as you wish.


