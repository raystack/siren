import Tabs from "@theme/Tabs";
import TabItem from "@theme/TabItem";
import CodeBlock from '@theme/CodeBlock';
import siteConfig from '/docusaurus.config.js';

# Slack
export const apiVersion = siteConfig.customFields.apiVersion
export const defaultHost = siteConfig.customFields.defaultHost

|||
|---|---|
|**type**|`slack`|


A Slack receiver in Siren tied to a Slack workspace. The implementation uses [slack-go](https://github.com/slack-go/slack) library to send notification to a [Slack app](#initializing-a-slack-app). This Slack app must be installed in the required workspaces and added to the required channels. Siren helps with the installation flow by automating the exchanging code for access token [flow](https://api.slack.com/legacy/oauth#authenticating-users-with-oauth__the-oauth-flow).


## Initializing a Slack App
Here is the list of actions one need to take to attach a Slack app to Siren. 

1. Create a Slack app and configure these permissions. Visit [this](https://api.slack.com/apps). If you already have an app, make sure permissions mentioned below are there.
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
4. Setup a redirection server. You can use `localhost` as well. This must be a `https` server. Slack will call this server once we install the app in any workspace.
5. Install your app to a workspace. Visit `Manage Distribution` section on the App Dashboard. Click the `Add to Slack` Button.
6. This will prompt you to the OAuth Consent screen. Make sure you have selected the correct Slack Workspace by verifying the dropdown in the top-right corner. Click `Allow`.
7. Copy the `code` that you received from Slack redirection URL query params and use this as `auth_code` inside `create receiver` payload.

In order to send Slack notifications via Siren API, you need to create its receiver.   To add a new slack receiver you need all credentials mentioned in this [configuration](#configurations-in-api).


## Configurations in API

```json
"configurations": {
    "client_id": <string>,
    "client_secret": <string>,
    "auth_code": <string>
}
```

## Configurations Stored in DB

```json
"configurations": {
    "token": <encrypted string>
}
```

Creating a slack receiver involves exchanging the auth code for a token with slack oAuth server. Siren will need the auth code, client id, client secret and optional label metadata.

**Example**

<Tabs groupId="api">
  <TabItem value="cli" label="CLI" default>

```bash
$ siren receiver create --file slack_receiver.yaml
```

  </TabItem>
  <TabItem value="http" label="HTTP">
    <CodeBlock className="language-bash">
    {`$ curl --request POST
  --url `}{defaultHost}{`/`}{apiVersion}{`/receivers
  --header 'content-type: application/json'
  --data-raw '{
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
}'`}
    </CodeBlock>
  </TabItem>
</Tabs>

On success, this will store the encrypted app token for that particular slack workspace and use it for sending out notifications.


## Subscription

Slack has additional `SubscriptionConfig` where user could routes the notification to a specific channel or individual in a workspace. Here is the subscription config.

```json
"configurations": {
    "channel_name": <string>,
    "channel_type": <string>
}
```

The `channel_type` has two possible enum values `channel` and `user`. The default value of this is `channel`. If `channel_type` is channel, `channel_name` should be a slack channel handle in the workspace e.g. `#odpf-critical`. If one wants to send notification to an individual, `channel_type` needs to be `user` and `channel_name` needs to be populated with slack user's e-mail e.g. `user1@odpf.io`.


## Message Payload

### Contract

Payload format of slack needs to follow [slack chat.postMessage API](https://api.slack.com/methods/chat.postMessage) contract.

```yaml
channel: <string>
text: <string>
username: <string>
icon_emoji: <string>
icon_url: <string>
link_names: <boolean>
attachments:
  - <key1>: <any>
    <key2>: <any>
  - <key3>: <any>
    <key4>: <any>
    .
    .
```

### Default Alert Template

Siren has a slack default notification [template](../../../plugins/receivers/slack/config/default_alert_template_body.goyaml) used by all alert notifications.
