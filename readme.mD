#### Four unique ways to deploy and scale your containerized app on Google Cloud.

I gave this talk at #DevFestNairobi22.

The slides for the talk are [here](https://docs.google.com/presentation/d/16NU_pAGr0hmur5mnosQ4iOtK7z3wEvDG9e7wfAJFaUE/edit?usp=sharing)

You can follow these docs to run the containerized demo app on `Cloud Run`, `App Engine`, `GKE` and `Compute Engine`

The `config` directory contains the scripts and manifests. 

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

My project ID is: `devfest-demos-364710` which is what I use in the following demo commands. You will need to replace it with your project's id.

#####  Running the app in your local machine
You can run the app using `go run server/main.go`
The app will be running on `http://0.0.0.0:8080/`


Once you verify the app is working, it's now time to package it in Docker using these commands.
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

With our app packaged in docker and pushed to the artifact registry, we can run it in cloud run using below configuration. The config for running it on cloud run is in the file `config/cloudbuild.yaml`.

We only have one build step since we already built and pushed the image.
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

To allow internet traffic so you access the app;
 - Select the just created cloud run service in your GCP console. 
 - In the new principals textbox, enter the value `allUsers`
 - From the Role dropdown menu, select the Cloud Run Invoker role.
 - Click Save.

 You can now access the app running on cloud run with the url that provided by cloud run.

##### Deploying App Engine Flex
We can also deploy our app to App engine flex environment where it will be built as docker image.

The App Engine Flex settings are in the file `app.yaml`. 

We set the runtime as custom. This means App engine will the Dockerfile to build the image using [cloud build](https://cloud.google.com/build) and deploy it to flex environment..
```
runtime: custom
env: flex
resources:
  cpu: 1
  memory_gb: 0.5
  disk_size_gb: 10
liveness_check:
  path: "/ping"
  check_interval_sec: 30
  timeout_sec: 4
  failure_threshold: 2
  success_threshold: 2
automatic_scaling:
  min_num_instances: 1
  max_num_instances: 3
  cool_down_period_sec: 180
  cpu_utilization:
    target_utilization: 0.9
  target_concurrent_requests: 100
```
The command to deploy the app to flex is `gcloud app deploy`.  After the deploment you will see the url of the app on the console. 

You can also access the app through `gcloud app browse`

##### Deploying to GKE
To deploy the same containerised app to GKE, you first need to [create a cluster](https://cloud.google.com/kubernetes-engine/docs/deploy-app-cluster) and fetch the credentials using these commands. 
```
#Create a cluster
gcloud container clusters create defvest-demo-cluster-1 --zone=europe-west1

#Fetch credentials for the cluster
gcloud container clusters get-credentials defvest-demo-cluster-1 --zone=europe-west1

```
We will use the Kubernetes manifests in `config/k8s` directory. i.e 
- `deployment.yaml` 
- `service.yaml` 
- `hpa.yaml` 

To apply the manifests to your new GKE cluster, run

`kubectl apply -f config/k8s`

To get the Public IP of the service so you can access the application run
`kubectl get svc`

##### Deploying to Container Optimized VM

We will also deploy the same image on a Container Optimized VM on Compute Engine. The command for create a VM with our image  is in `config/compute/run.sh`

```
gcloud compute instances \
create-with-container \
devfest-containers-demo \
--machine-type g1-small \
--zone=europe-west1-b \
--container-image europe-west1-docker.pkg.dev/devfest-demos-364710/gcp-containers-demo/demoapp:latest

```
You need to allow http traffic on the [firewall](https://cloud.google.com/vpc/docs/firewalls) and then you can access the app running on the vm via the attached ip.
Remember to clean up your resources after to avoid additional costs. 