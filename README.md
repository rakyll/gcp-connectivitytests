# Google Cloud Connectivity Tests

## Installation

Linux 64-bit:

```
$ curl http://storage.googleapis.com/jbd-releases/gcp-connectivitytests-0.0.1-linuxamd64 > gcp-connectivitytests && chmod +x gcp-connectivitytests
```

Use Google Application Default Credentials to authenticate.

----

See the following sections to generate boilerplate manifests
to run tests as a post-commit hook for your master branch.

## Travis

Download a JSON service account key from https://console.cloud.google.com/iam-admin/serviceaccounts.
You need to have Travis CLI installed in order to use this command.

Run the following command to generate the Travis CI manifest:

```
$ gcp-connectivitytests -gen=travis -project=PROJECT_ID -secretkey=key.json > .travis.yml
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

## Circle CI

Download a JSON service account key from https://console.cloud.google.com/iam-admin/serviceaccounts.
Provide the contents of the key as the GCLOUD_SERVICE_KEY env variable.

Run the following command to generate the Circle CI manifest:

```
$ gcp-connectivitytests -gen=circleci -project=PROJECT_ID > .circleci/config.yml
version: 2
jobs:
  build:
    docker:
      - image: google/cloud-sdk

    steps:
      - run: |
          apt-get install wget -y
          wget https://storage.googleapis.com/jbd-releases/gcp-connectivitytests-0.0.1-linuxamd64 && chmod +x ./gcp-connectivitytests-0.0.1-linuxamd64
          echo $GCLOUD_SERVICE_KEY > key.json
          GOOGLE_APPLICATION_CREDENTIALS=key.json ./gcp-connectivitytests-0.0.1-linuxamd64 -project=PROJECT_ID
```

See [rakyll/gcp-connectivitytests-example](https://github.com/rakyll/gcp-connectivitytests-example) as an
example repo setup with Travis and Circle CI.