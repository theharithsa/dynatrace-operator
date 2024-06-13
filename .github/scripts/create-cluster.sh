#!/usr/bin/env bash

set -x

DEFAULT_TIMEOUT="60m" # K8s takes <10min, Openshift >40min
desired_state="environment-deployed"

kubectl version

echo "Creating environment '$FLC_ENVIRONMENT' in namespace '$FLC_NAMESPACE'"

kubectl get flcenvironments --namespace $FLC_NAMESPACE

echo "Patching environment '$FLC_ENVIRONMENT' to 'deployed'"
kubectl patch --namespace $FLC_NAMESPACE --type merge --patch '{"spec": {"desiredState": "environment-deployed"}}' flcenvironment $FLC_ENVIRONMENT

echo "Waiting up to '$DEFAULT_TIMEOUT' for successful deployment of environment '$FLC_ENVIRONMENT'"
kubectl wait --namespace $FLC_NAMESPACE --timeout="$DEFAULT_TIMEOUT" --for=condition=InTransition=false flcenvironment $FLC_ENVIRONMENT

if [ "$PLATFORM" == "ocp" ]; then
  echo "PLATFORM is set to 'ocp'. Waiting for 5 minutes..."
  sleep 300  # Wait for 5 minutes (300 seconds)
  echo "done"
fi

echo "Checking currentState='$desired_state' for '$FLCENV_NAME'..."
flc_state=$(kubectl get flcenvironment "$FLCENV_NAME" --namespace "$NAMESPACE" -ojsonpath="{.status.currentState}")
if [[ "$flc_state" != "$desired_state" ]]; then
  echo "Pipeline deployment did not reach expected state '$desired_state', currentState: ${flc_state}"
  exit 1
fi
echo "done"
