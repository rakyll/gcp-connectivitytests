# Google Cloud Connectivity Tests

## Installation

Linux 64-bit:

```
$ curl http://storage.googleapis.com/jbd-releases/gcp-connectivitytests-0.0.1 > gcp-connectivitytests && chmod +x gcp-connectivitytests
```

Use Google Application Default Credentials to authenticate.

See the following sections to generate boilerplate CI manifests
to run tests as a post-commit hook for your master branch.

## Travis

Download a service account key from https://console.cloud.google.com/iam-admin/serviceaccounts.
You need to have Travis CLI installed in order to use this command.

```
$ gcp-connectivitytests -project=PROJECT_ID -gen=travis -secretkey=path-to-secret-key.json
```
