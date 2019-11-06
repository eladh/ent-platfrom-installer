#!/bin/sh

echo " -------------"
echo "| GCLOUD INFO |"
echo " -------------"
gcloud info

echo " ------------------------"
echo "| KUBECTL CLIENT VERSION |"
echo " ------------------------"
kubectl version --client
echo ""

echo " ---------------------"
echo "| HELM CLIENT VERSION |"
echo " ---------------------"
helm version --client
echo ""

echo " ----------------"
echo "| AUTHENTICATE  |"
echo " ----------------"
export GOOGLE_APPLICATION_CREDENTIALS=/home/appuser/license.json
gcloud auth activate-service-account --key-file=/home/appuser/license.json
echo ""

echo " ----------------"
echo "| RUN           |"
echo " ----------------"
#  ---- Use to debug the agent container ---
# tail -f /dev/null
/home/appuser/platformInstaller --config ./setup.yaml ${COMMAND}
echo ""

