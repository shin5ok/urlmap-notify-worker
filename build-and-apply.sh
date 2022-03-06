#!/bin/bash
export TAG=$(date '+%Y%m%d%H%M')
if gcloud beta builds submit --pack=image=us-central1-docker.pkg.dev/$PROJECT/myrepo/urlmap-notify-worker:$TAG;
then
    envsubst < manifest.yaml | kubectl apply -f -
fi
