#!/bin/bash

# delete webhook config
kubectl delete mutatingwebhookconfiguration label-hook-cfg

# delete webhook server deployment
kubectl delete deployment label-hook-deployment -n label-hook

# rebuild & push to hub
CURRENT_DIR=$(cd "$(dirname "$0")";pwd)
sh $CURRENT_DIR/build_image.sh

# file CABundle
CA_BUNDLE=$(kubectl config view --raw --flatten -o json | jq -r '.clusters[] | .cluster."certificate-authority-data"')
sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g" $CURRENT_DIR/../resources/mutating-admission-webhook.yaml > $CURRENT_DIR/../resources/mutating-admission-webhook-tmp.yaml

# deploy webhook server & webhook config
kubectl apply -f $CURRENT_DIR/../resources/rbac.yaml
kubectl apply -f $CURRENT_DIR/../resources/deployment.yaml
kubectl apply -f $CURRENT_DIR/../resources/service.yaml
kubectl apply -f $CURRENT_DIR/../resources/mutating-admission-webhook-tmp.yaml
