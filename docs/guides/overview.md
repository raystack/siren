# Usage

Siren comes with built in command which perform specific tasks.

## List of available commands

1. Serve
    - Runs the Server  `$ go run main.go serve`

2. Migrate
    - Runs the DB Migrations `$ go run main.go migrate`

3. Upload
    - Parses a YAML File in specified format to upsert templates and rules(
      alerts) `$ go run main.go upload fileName.yaml`. Read more about the Rules and Templates [here](../concepts).
