# nemesis

Nemesis is a tool for auditing platform configurations for measuring compliance. It is meant as a read-only view into Cloud platforms such that necessary audits can be performed, which can result in actions to take.

## Usage
You can install `nemesis` as a binary on your machine, or run it as a docker container.

The following line demonstrates basic usage to invoke `nemesis` and output results into your terminal. This assumes that you have valid GCP credentials on the host you are running on:
```
nemesis --project.filter="my-project" --reports.stdout.enable
```

You can utilize a service account credential file to perform `nemesis` runs as the service account user:
```
# Set the environment variable that the Google Auth library expects
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json

# Now run nemesis as the service account
nemesis --project.filter="my-awesome-project" --reports.stdout.enable
```

You can also combine `nemesis` with tools like `jq` to do parsing or formatting a JSON file:
```
# Output a nemesis report to a local JSON file that is formatted in a readable way

nemesis --project.filter="my-project" --reports.stdout.enable | jq . >> report.json
```

You can scan multiple projects in a single `nemesis` run using a simple regular expression:
```
nemesis --project.filter="my-business-unit-projects-*"
```

`nemesis` reports can be directly shipped to a GCP Pub/Sub topic for direct ingestion into another system:
```
nemesis --project.filter="my-project" --reports.pubsub.enable --reports.pubsub.project="my-reporting-project" --reports.pubsub.topic="nemesis-reports'
```

All flags for `nemesis` have an equivalent environment variable you can use for configuration. The table under Flags indicates the equivalencies:
```
# Configure many settings before running
export NEMESIS_METRICS_ENABLED="true"
export NEMESIS_METRICS_GATEWAY="prometheus-pushgateway.example.com:9091"
export NEMESIS_PROJECT_FILTER="my-project"
export NEMESIS_ONLY_FAILURES="true"
export NEMESIS_ENABLE_STDOUT="true"

# Now run the scan
nemesis
```


## Flags
`nemesis` has a number of flags that can be invoked either using the command line flag or the equivalent environment variable. The following table describes their usage:

| Flag | Environment Variable | Required | Description | Example Flag Usage |
|------|----------------------|----------|-------------|--------------------|
| project.filter                        | `NEMESIS_PROJECT_FILTER`              | yes   | (String) The project filter to perform audits on                                          | `--project.filter="my-project"`   |
| compute.instance.allow-ip-forwarding  | `NEMESIS_COMPUTE_ALLOW_IP_FORWARDING` | no    | (Bool) Indicate whether instances should be allowed to perform IP forwarding              | `--compute.instance.allow-ip-forwarding`              |
| compute.instance.allow-nat            | `NEMESIS_COMPUTE_ALLOW_NAT`           | no    | (Bool) Indicate whether instances should be allowed to have external (NAT) IP addresses   | `--compute.instance.allow-nat`                        |
| compute.instance.num-interfaces       | `NEMESIS_COMPUTE_NUM_NICS`            | no    | (String) The number of network interfaces (NIC) that an instance should have (default 1)  | `--compute.instance.num-interfaces=1`                 |
| container.oauth-scopes                | `NEMESIS_CONTAINER_OAUTHSCOPES `      | no    | (String) A comma-seperated list of OAuth scopes to allow for GKE clusters (default <br>"https://www.googleapis.com/auth/devstorage.read_only,<br>https://www.googleapis.com/auth/logging.write,<br>https://www.googleapis.com/auth/monitoring,<br>https://www.googleapis.com/auth/servicecontrol,<br>https://www.googleapis.com/auth/service.management.readonly,<br>https://www.googleapis.com/auth/trace.append") | `--container.oauth-scopes="..."` |
| iam.sa-key-expiration-time            | `NEMESIS_IAM_SA_KEY_EXPIRATION_TIME`  | no    | (String) The time in days to allow service account keys to live before being rotated (default "90") | `--iam.sa-key-expiration-time="90"` |
| iam.user-domains                      | `NEMESIS_IAM_USERDOMAINS`             | no    | (String) A comma-separated list of domains to allow users from                            | `--iam.user-domains="google.com"` |
| metrics.enabled                       | `NEMESIS_METRICS_ENABLED`             | no    | (Boolean) Enable Prometheus metrics                                                       | `--metrics.enabled` |
| metrics.gateway                       | `NEMESIS_METRICS_GATEWAY`             | no    | (String) Prometheus metrics Push Gateway (default "127.0.0.1:9091")                       | `--metrics.gateway="10.0.160.12:9091"` |
| reports.only-failures                 | `NEMESIS_ONLY_FAILURES`               | no    | (Boolean) Limit output of controls to only failed controls                                | `--reports.only-failures` |
| reports.stdout.enable                 | `NEMESIS_ENABLE_STDOUT`               | no    | (Boolean) Enable outputting report via stdout                                             | `--reports.stdout.enable` |
| reports.pubsub.enable                 | `NEMESIS_ENABLE_PUBSUB`               | no    | (Boolean) Enable outputting report via Google Pub/Sub                                     | `--reports.pubsub.enable` |
| reports.pubsub.project                | `NEMESIS_PUBSUB_PROJECT`              | no    | (Boolean) Indicate which GCP project to output Pub/Sub reports to                         | `--reports.pubsub.project="my-project"` |
| reports.pubsub.topic                  | `NEMESIS_PUBSUB_TOPIC`                | no    | (Boolean) Indicate which topic to output Pub/Sub reports to (default "nemesis")           | `--reports.pubsub.topic="nemesis-reports"` |

## Motivation

`nemesis` was created out of a need to generate compliance and auditing reports quickly and in consumable formats. This tool helps audit against GCP security standards and best practices. We implement, as a baseline security metric:
* [CIS Controls for GCP](https://www.cisecurity.org/benchmark/google_cloud_computing_platform/)
* [GKE Hardening Guidelines](https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster)
* [Default Project Metadata](https://cloud.google.com/compute/docs/storing-retrieving-metadata#default)

We strive to encourage best practices in our environment for the following GCP services:
* [Identity and Access Management (IAM)](https://cloud.google.com/iam/docs/using-iam-securely)
* [Google Cloud Storage (GCS)](https://cloud.google.com/storage/docs/access-control/using-iam-permissions)
* [Google Compute Engine (GCE)](https://cloud.google.com/compute/docs/access/)
* [Google Kubernetes Engine (GKE)](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-admin-overview#configuring_cluster_security)

## Maintainers

@TaylorMutch

## Contributions

See Contributions