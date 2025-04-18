gcloud run deploy go-http-function \
    --source . \
    --function ComposeHaiku \
    --base-image go123 \
    --region europe-north2 \
    --allow-unauthenticated
