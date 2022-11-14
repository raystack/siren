# Job

Job is a task that only runs once and then the process is terminated. Job could be scheduled with Cron or triggered manually. Siren is currently only having a one very-specific job.

## Queue Cleanup Job (Postgres Queue Only)

This job requires the same config like server. This job cleans up all stored `published` messages in queue with last updated more than specific threshold (default 7 days) from `now()` and optionally cleaning up all `pending` messages in queue with last updated more than specific threshold (default 7 days) from `now()`.

Run this command to clean up all published messages with age 7 days and more.

```bash
$ siren job run cleanup_queue --config config.yaml
```

Run this command to clean up all published messages with age 7 days and more and clean up all pending messages with age 14 days and more.

```bash
$ siren job run cleanup_queue --pending 336h --config config.yaml
```