# Google Cloud Connectivity Tests

## Installation

Linux 64-bit:

```
$ curl http://storage.googleapis.com/jbd-releases/gcp-connectivitytests-0.0.1-linuxamd64 > gcp-connectivitytests && chmod +x gcp-connectivitytests
```

If you have Go installed, use:

Use Google Application Default Credentials to authenticate.

----

See the following sections to generate boilerplate manifests
to run tests as a post-commit hook for your master branch.

## Travis

Download a JSON service account key from https://console.cloud.google.com/iam-admin/serviceaccounts.
You need to have Travis CLI installed in order to use this command.

```
$ gcp-connectivitytests -gen=travis -project=PROJECT_ID -secretkey=key.json
branches:
  only:
    - master

before_install:
 - openssl aes-256-cbc -K $encrypted_e53f4ef918d8_key -iv $encrypted_e53f4ef918d8_iv -in key.json.enc -out key.json -d

install:
 - wget https://storage.googleapis.com/jbd-releases/gcp-connectivitytests-0.0.1-linuxamd64 && chmod +x ./gcp-connectivitytests-0.0.1

script:
 - GOOGLE_APPLICATION_CREDENTIALS=key.json ./gcp-connectivitytests-0.0.1 -project=PROJECT_ID
```

# Circle CI

Coming soon.