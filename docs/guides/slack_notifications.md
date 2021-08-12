# Slack Notification

Siren can send Slack messages to channels or users via an HTTP API interface. It uses a preconfigured Slack app with
required permissions. The abstractions this feature provides is that channel name/user email can be used to send out
Slack messages and like all other features in Siren, this feature is also multi-tenant, which means you can choose the
workspace to send out the messages in the request payload.

In order to send slack notifications via Siren Apis, you need to attach a slack App to siren. This slack app then must
be installed in the required workspaces and added to the required channels.

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
7. Copy the `code` that you received from Slack redirection URL query params.

## Exchanging Access Token for Auth Code

Now siren can be used to exchange the `code` for an `access_token` native to a workspace. This access_token will be
stored in encrypted format by Siren and will be used to send out Slack message in that workspace using the Slack App
created above.

Let's see an example.

Please note that, in order for this exchange to be successful, you will need to mention, the `client_id` and
`client_secret` of the Slack App, inside siren's config file along with a 32 character encryption key.

```yaml
slack_app:
   client_id: '123.456'
   client_secret: 'secret'
encryption_key: 'some_random_32_character_string'
```

```text
curl -X 'POST' \
  'http://localhost:3000/oauth/slack/token' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "code": "auth-code",
  "workspace": "workspace_name"
}'
```

If this request succeeds Siren will store the token in DB in encrypted format. In the next steps, you can see how it
will be used to send out Slack messages. Otherwise, logs will help you with the exact error message.

### Send Slack Notification

You need to specify which entity(workspace name as mentioned above), receiver name, type(user/channel) along with
message body(either in the Markdown format or with [blocks](https://api.slack.com/block-kit))

```text
curl --location --request POST 'localhost:3000/notifications?provider=slack' \
--header 'Content-Type: application/json' \
--data-raw '{
    "entity": "odpf",
    "receiver_name": "abhishek.sah@odpf.io",
    "receiver_type": "user",
    "message": "hello world",
    "blocks": [
        {
            "type": "section",
            "text": {
                "type": "mrkdwn",
                "text": "Hi there"
            }
        },
        {
            "type": "section",
            "fields": [
                {
                    "type": "mrkdwn",
                    "text": "*Type:*\nComputer (laptop)"
                },
                {
                    "type": "mrkdwn",
                    "text": "*When:*\nSubmitted Aut 10"
                }
            ]
        },
        {
            "type": "actions",
            "elements": [
                {
                    "type": "button",
                    "text": {
                        "type": "plain_text",
                        "emoji": true,
                        "text": "Approve"
                    },
                    "style": "primary",
                    "value": "click_me_123"
                },
                {
                    "type": "button",
                    "text": {
                        "type": "plain_text",
                        "emoji": true,
                        "text": "Deny"
                    },
                    "style": "danger",
                    "value": "click_me_123"
                }
            ]
        }
    ]
}
```

### Terminology of the request body

| Term          | Description                                                                     |
|---------------|---------------------------------------------------------------------------------|
| Entity        | slack workspace name. Assumption: The Slack App was installed in this workspace |
| receiver_type | channel/user                                                                    |
| receiver_name | email of user/name of channel                                                   |
| message       | markdown format message body                                                    |
| blocks        | block array [see](https://api.slack.com/block-kit). It overrides message!       |

**Note:**

1. Either of block or message is required.
2. Blocks overrides message.
3. For simple text messages, you can use `message` key.
