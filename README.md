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
branches:
  only:
    - master

before_install:
 - openssl aes-256-cbc -K $encrypted_e53f4ef918d8_key -iv $encrypted_e53f4ef918d8_iv -in path-to-secret-key.json.enc -out secret-key.json -d

install:
 - wget https://storage.googleapis.com/jbd-releases/gcp-connectivitytests-0.0.1 && chmod +x ./gcp-connectivitytests-0.0.1

# Build the website
script:
 - GOOGLE_APPLICATION_CREDENTIALS=secret-key.json ./gcp-connectivitytests-0.0.1 -project=PROJECT_ID
```
