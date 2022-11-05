gcloud artifacts repositories create gcp-containers-demo --repository-format=docker \
--location=europe-west1 --description="DevFest Docker repository"

gcloud auth configure-docker europe-docker.pkg.dev


#Build the image and tag it using our artifact registry and project name
docker build -t europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest .

#Push the docker image to artifact registry
docker push europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest