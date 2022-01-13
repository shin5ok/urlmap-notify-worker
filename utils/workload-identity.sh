if ! ( kubectl get ns | grep urlmap > /dev/null );
then
    kubectl create namespace urlmap
fi
kubectl create serviceaccount --namespace urlmap urlmap-notify-worker
gcloud iam service-accounts create urlmap-notify-worker
gcloud projects add-iam-policy-binding $PROJECT     --member "serviceAccount:urlmap-notify-worker@$PROJECT.iam.gserviceaccount.com"     --role "roles/secretmanager.secretAccessor"
gcloud projects add-iam-policy-binding $PROJECT     --member "serviceAccount:urlmap-notify-worker@$PROJECT.iam.gserviceaccount.com"     --role "roles/pubsub.subscriber"
gcloud iam service-accounts add-iam-policy-binding     --role roles/iam.workloadIdentityUser     --member "serviceAccount:$PROJECT.svc.id.goog[urlmap/urlmap-notify-worker]"     urlmap-notify-worker@$PROJECT.iam.gserviceaccount.com
kubectl annotate serviceaccount \
    --namespace urlmap urlmap-notify-worker \
    iam.gke.io/gcp-service-account=urlmap-notify-worker@$PROJECT.iam.gserviceaccount.com
