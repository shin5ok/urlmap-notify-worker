#!/bin/bash
TAG=${TAG:-0.01}
gcloud beta builds submit --pack=image=us-central1-docker.pkg.dev/shingo-ar-proto/myrepo/urlmap-notify-worker:$TAG