#!/usr/bin/env bash
set -euo pipefail

gcloud run deploy go-http-function \
    --source . \
    --function ComposeHaiku \
    --base-image go123 \
    --region europe-north2 \
    --env-vars-file=env.yaml \
    --allow-unauthenticated
