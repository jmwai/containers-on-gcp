gcloud compute instances \
create-with-container \
devfest-containers-demo \
--machine-type g1-small \
--zone=europe-west1-b \
--container-image europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest