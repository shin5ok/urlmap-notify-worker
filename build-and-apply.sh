#!/bin/bash
export TAG=$(date '+%Y%m%d%H%M')
gcloud beta builds submit --pack=image=us-central1-docker.pkg.dev/shingo-ar-proto/myrepo/urlmap-notify-worker:$TAG
envsubst < manifest.yaml | kubectl apply -f -
