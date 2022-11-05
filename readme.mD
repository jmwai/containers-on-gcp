### GCP Containers Demo.

This projects shows how to run a containerized app on GCP.

The slides for the talk are [here](https://docs.google.com/presentation/d/18Z38kkWKpe-K0IAQjKBSr2QI07DKg0krQ8yfK6hqY5c/edit?usp=sharing)

The config/ directory contains scripts and yaml files that you can execute to deploy the app on GCP.
 
The config directory has 
- cloud_run
- compute
- flex
- k8s
- registry

#### Prerequisites 
- Create a Google Cloud account [here](https://cloud.google.com/) if you don't have one.
- Install the Google Cloud CLI [here](https://cloud.google.com/sdk/docs/install)
- You will also need to create a [Project](https://cloud.google.com/resource-manager/docs/creating-managing-projects) and enable these APIs
    
    - Cloud Run
    - Artifact Registry
    - App Engine
    - GKE
    - Compute Engine
- We will use Artifact Registry to store the docker image. You can learn about and create an Artifact Registry for your project [here](https://cloud.google.com/artifact-registry)

The demo app use Google Maps API, so to run it you need to get an API Key [here.](https://developers.google.com/maps/documentation/embed/get-api-key) 

Add the key to the .env file `GOOGLE_MAPS_API_KEY="your-api-key"`

My project ID is: `devfest-demos-364710` which is what I use in the following demo commands

#####  Running the app in your local machine
You can run the app using `go run server/main.go`
The app will be runnging on `http://0.0.0.0:8080/`


Once you verify the app is working it's now time to package it in Docker using these commands.
##### Build the Image
```

#Build the image and tag it using our artifact registry and project name
docker build -t europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest .
```
My artifact registry is in `europe-west1` and my repositry name is `gcp-containers-demo`
Learn more about naming convention for artifact registry [here](https://cloud.google.com/artifact-registry/docs/docker/pushing-and-pulling) 

##### Run the app in docker on your machine

```
docker run -p 8080:8080 europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest
```

 The app will be runnging on `http://0.0.0.0:8080/`
##### Pushing the image to Artifact Registry
After verifying the app is running locally we need to push the image to the artifact registry so we can access it from the cloud services we will be testing.

`docker push europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest`

##### Deploying to Cloud Run

With our app packaged in docker and pushed to the artifact registry, we can run it in cloud run using below configuration. The config for running it on Cloud is in the file `config/cloudbuild.yaml` with this yaml code.
```
steps:
# Deploy existing container image to Cloud Run
- name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
  entrypoint: gcloud
  args: 
  - 'run' 
  - 'deploy'
  - 'gcp-containers-demo'
  - '--image'
  - 'europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest'
  - '--region'
  - 'europe-west1'
```

To deploy run `gcloud builds submit` in the same directoy where `cloudbuild.yaml` is located.

After the deploment you will see the url of the app running on cloud run on the console.

You need to allow internet traffic to access the app.

Select the service in the console. 
 - In the New principals textbox, enter the value allUsers
 - From the Role dropdown menu, select the Cloud Run Invoker role.

Click Save.

##### Deploying App Engine Flex
We can also deploy our app to App engine flex where it will be built as docker image.

The Flex configs are in the file `app.yaml`. 

We set the runtime as custom. This means cloud build will use our a Dockerfile to build and deploy the app to App Engine Flex.
```
runtime: custom 
env: flex
manual_scaling:
  instances: 1
resources:
  cpu: 1
  memory_gb: 0.5
  disk_size_gb: 10
```
The command to deploy the app to flex is `gcloud app deploy`. App engine uses the Dockerfile to build the image using [cloud build](https://cloud.google.com/build) and deploying it to flex environment. After the deploment you will see the url of the app on the console. 

You can also access the app through `gcloud app browse`

##### Deploying to GKE
To deploy the same containerised app to GKE, you first need to create a cluster and fetch the credentials using these commands. 
```
#Create a cluster
gcloud container clusters create defvest-demo-cluster-1 --zone=europe-west1

#Fetch credentials for the cluster
gcloud container clusters get-credentials defvest-demo-cluster-1 --zone=europe-west1

```
We will use the configs in `config/k8s` directory. i.e 
- `deployment.yaml` 
- `service.yaml` 
- `hpa.yaml` 

This command will apply the configs to our GKE cluster

`kubectl apply -f config/k8s`

To get the IP so you can access the application run
`kubectl get svc`

##### Deploying to Container Optimized VM

We will can also deploy the same image in a Container Optimized VM on Compute Engine. The command for this is in `config/compute/run.sh`

```
gcloud compute instances \
create-with-container \
gcp-containers-demo \
--container-image europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest

```