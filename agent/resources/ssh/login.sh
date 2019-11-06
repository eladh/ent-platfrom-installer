#!/usr/bin/env bash

CLUSTER_NAME=$1

gcloud auth activate-service-account --key-file=/tmp/license.json
gcloud container clusters get-credentials ${CLUSTER_NAME} --zone us-central1-a --project devops-consulting

helm init --client-only
helm plugin install https://github.com/rimusz/helm-tiller

helm repo add jfrog https://charts.jfrog.io/
helm repo update

kubectl get svc