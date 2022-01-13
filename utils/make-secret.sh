#!/bin/bash
echo -n $SECRET | gcloud secrets create SLACK_URL --data-file=-