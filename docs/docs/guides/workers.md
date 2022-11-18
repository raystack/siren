# Workers

Siren has a notification features that utilizes queue to publish notification messages. More concept about notification could be found [here](../concepts/notification.md). The architecture requires a detached worker running asynchronously and polling queue periodically to dequeue notification messages and publish them. By default, Siren server run this asynchronous worker inside it. However it is also possible to run the worker as a different process. Currently there are two possible workers to run
1. **Notification message handler:** this worker periodically poll and dequeue messages from queue, process the messages, and then publish notification messages to the specified receivers. If there is a failure, this handler enqueues the failed messages to the dlq.
1. **Notification dlq handler:** this worker periodically poll and dequeue messages from dead-letter-queue, process the messages, and then publish notification messages to the specified receivers. If there is a failure, this handler enqueues the failed messages back to the dlq.



## Running Workers as a Different Process

Siren has a command to start workers. Workers use the same config like server does.

```bash
Command to start a siren worker.

Usage
  siren worker start <command> [flags]

Core commands
  notification_dlq_handler    A notification dlq handler
  notification_handler        A notification handler

Inherited flags
  --help   Show help for command

Examples
  $ siren worker start notification_handler
  $ siren server start notification_handler -c ./config.yaml
```

Starting up a worker could be done by executing.

```bash
$ siren worker start notification_handler
```