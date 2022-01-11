# Receivers

Receivers represent a notification medium, which can be used to define routing configuration in the monitoring
providers, to control the behaviour of how your alerts are notified. Few examples: Slack receiver, HTTP receiver,
Pagerduty receivers etc. Currently, Siren supports these 3 types of receivers. Configuration of each receiver depends on
the type.

You can use receivers to send notifications on demand as well as on certain matching conditions. Subscriptions use
receivers to define routing configuration in monitoring providers. For eg. Cortex-metrics uses alertmanager for routing
alerts. With Siren subscriptions, you will be able to manage routing in Alertmanager using pre-registered receivers.

## API Interface

### Create a Receiver

**Type: Slack**

Using a slack receiver you will be able to send out Slack notification using its send API. You can also use it to route
alerts using Subscriptions whenever an alert matches the conditions of your choice. Check the required permissions of
the Slack App [below](#permissions-and-auth-settings-for-slack-receivers).

Creating a slack receiver involves exchanging the auth code for a token with slack oauth server. Siren will need the
auth code, client id, client secret and optional label metadata.

Example

```text
POST /v1beta1/receivers HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 228

{
    "name": "doc-slack-receiver",
    "type": "slack",
    "labels": {
        "team": "siren-devs"
    },
    "configurations": {
        "client_id": "abcd",
        "client_secret": "xyz",
        "auth_code": "123"
    }
}
```

On success, this will store the app token for that particular slack workspace and use it for sending out notifications.

**Type: Pagerduty**

```text
POST /v1beta1/receivers HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 182

{
    "name": "doc-pagerduty-receiver",
    "type": "http",
    "labels": {
        "team": "siren-devs"
    },
    "configurations": {
        "url": "http://localhost:4000"
    }
}
```

**Type: HTTP**

```text
POST /v1beta1/receivers HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 177

{
    "name": "doc-http-receiver",
    "type": "http",
    "labels": {
        "team": "siren-devs"
    },
    "configurations": {
        "url": "http://localhost:4000"
    }
}
```

### Update a Receiver

**Note:** While updating a receiver, you will have to make sure all subscriptions that are using this receivers get
refreshed(updated), since subscriptions use receivers to create routing configuration dynamically.

**Type: HTTP**

```text
PUT /v1beta1/receivers/61 HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 177

{
    "name": "doc-http-receiver",
    "type": "http",
    "labels": {
        "team": "siren-devs"
    },
    "configurations": {
        "url": "http://localhost:4001"
    }
}
```

### Getting a receiver

```text
GET /v1beta1/receivers/61 HTTP/1.1
Host: localhost:3000
```

### Getting all receivers

```text
GET /v1beta1/receivers HTTP/1.1
Host: localhost:3000
```

### Deleting a receiver

```text
DELETE /v1beta1/receivers/61 HTTP/1.1
Host: localhost:3000
```

### Sending a message/notification

The types that supports sending messages using API are:

**Type: Slack**

```text
GET /v1beta1/receivers/51/send HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Content-Length: 399

{
    "slack": {
        "receiverName": "siren-devs",
        "receiverType": "channel",
        "blocks": [
            {
                "type": "section",
                "text": {
                    "type": "mrkdwn",
                    "text": "New Paid Time Off request from <example.com|Fred Enriquez>\n\n<https://example.com|View request>"
                }
            }
        ]
    }
}
```

Here we are using slack builder kit to construct block of messages, to send the message to channel name #siren-devs.

## CLI Interface

```text
Receivers are the medium to send notification for which we intend to mange configuration.

Usage:
  siren receiver [command]

Aliases:
  receiver, receivers

Available Commands:
  create      Create a new receiver
  delete      Delete a receiver details
  edit        Edit a receiver
  list        List receivers
  send        Send a receiver notification
  view        View a receiver details

Flags:
  -h, --help   help for receiver

Use "siren receiver [command] --help" for more information about a command.
```

### Permissions and Auth Settings for Slack Receivers

In order to send Slack notifications via Siren Apis, you need to create its receiver. This Slack app then must be
installed in the required workspaces and added to the required channels.

Siren helps with the installation flow by automating the exchanging code for access token
flow. [Reference](https://api.slack.com/legacy/oauth#authenticating-users-with-oauth__the-oauth-flow).

Here is the list of actions one need to take to attach a Slack app to Siren.

1. Create a Slack app with these permissions. Visit [this](https://api.slack.com/apps). If you already have an app, make
   sure permissions mentioned below are there.
2. Configure these permissions in the app:
   ```text
    channels:read
    chat:write
    groups:read
    im:read
    team:read
    users:read
    users:read.email
    ```
3. Enable Distribution
4. Setup a redirection server. You can use localhost as well. This must be a https server. Slack will call this server
   once we install the app in any workspace.
5. Install your app to a workspace. Visit `Manage Distribution` section on the App Dashboard. Click the `Add to Slack`
   Button.
6. This will prompt you to the OAuth Consent screen. Make sure you have selected the correct Slack Workspace by
   verifying the dropdown in the top-right corner. Click Allow.
7. Copy the `code` that you received from Slack redirection URL query params and use this inside create receiver
   payload.
