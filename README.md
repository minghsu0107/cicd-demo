# CI/CD Demo
[![Build Status](https://drone.pdlab.csie.ntu.edu.tw/api/badges/minghsu0107/cicd-demo/status.svg)](https://drone.pdlab.csie.ntu.edu.tw/minghsu0107/cicd-demo)

This is a tutorial of building CI/CD pipeline with DroneCI and ArgoCD. In this demo, we use DroneCI for running tests, publishing new images, and update image tags in the manifest repository. We then use ArgoCD for continuous delivery, synchronizing application states in the Kubernetes cluster with manifests maintained in the Git repository. This way of doing Kubernetes cluster management and application delivery is kown as GitOps. By applying GitOps, we can maintain a 'source of truth' for both the application code and infrastructure, improving system reliability and efficiency for your team.

## Prerequisites
1. A Drone server
    - [Github installation](https://docs.drone.io/server/provider/github/)
2. A K8s cluster
    - [K3d](https://k3d.io)
    - [minikube](https://minikube.sigs.k8s.io/docs/start/)
    - [K0s](https://github.com/k0sproject/k0s)
3. ArgoCD deployment
    - [All-in-one installation](https://argo-cd.readthedocs.io/en/stable/getting_started/#1-install-argo-cd)
4. A Github account and a Dockerhub account

## DroneCI
### Setup
After you have connected your Github account with Drone, you can browse all your repositories on Drone dashboard. Next, clone this repo, activate it and navigate to `Repositories -> cicd-demo -> settings` to add the following secrets:

![](https://i.imgur.com/CXj88J8.png)

- `docker_username`: your Dockerhub account
- `docker_password`: your Dockerhub password
- `git_email`: your email
- `git_username`: your Github account
- `git_password`: your Github password

Finally, replace `minghsu0107` with your Github and Dockerhub account in `.drone.yml`. Now any push or pull request will trigger a Drone pipeline. You can check details via `your repo -> setting -> webhook` on Github.
### Local Development
For local development, you will not want to push every change to your repo just for testing whether `.drone.yml` works. Instead, you can use [Drone CLI](https://docs.drone.io/cli/install/) to execute pipeline locally.

Login to Drone:
```bash
export DRONE_SERVER=<drone-server-url>
export DRONE_TOKEN=<drone-token> # check token under dashboard -> user setting
drone info
```
For example, you can run step `test` only by executing the following script under the project root:
```
drone exec --include=<pipline-step-name>
```
## ArgoCD
Please clone [the application manifest repository](https://github.com/minghsu0107/cicd-demo-manifests) first. This repo holds the application manifests and will be synced with ArgoCD later. The manifests are maintained by [Kustomize](https://github.com/kubernetes-sigs/kustomize), which is supported by ArgoCD out-of-the-box. 

If your repository is set to private, you need to configure access credentials on ArgoCD. Otherwise you can skip this step and create new app directly.

Credentials can be configured using Argo CD CLI:
```bash
argocd repo add <repo-url> --username <username> --password <password>
```
Or you can configure via UI. Navigate to `Settings/Repositories`; click Connect Repo using HTTPS button and enter credentials:

![](https://i.imgur.com/UAyNkte.png)

You will see something like:

![](https://i.imgur.com/XaMezBA.png)

Create new app:

![](https://i.imgur.com/gOD9h1b.png)

![](https://i.imgur.com/8XlNtDL.png)

![](https://i.imgur.com/JK76lnT.png)

Remember to place the repository with your own repo.

Now we have finish all preparations, and it's time to let the magic happen. Navigate to `/applications` and click SYNC button on your app in order to synchronize the cluster state:

![](https://i.imgur.com/RVH5QtL.png)

You can click your app to view details:

![](https://i.imgur.com/pconXQR.png)

As we can see, ArgoCD automatically sync the application to our desired state specified in `production` base. It also shows how all resources roll out in the cluster. With ArgoCD, we can not only have complete control over the entire application deployment but also track updates to branches, tags, or pinned to a specific version of manifests at a Git commit.
## Reference
- https://www.weave.works/technologies/gitops/
- https://argo-cd.readthedocs.io/en/stable/
- https://docs.drone.io
- https://hub.docker.com/r/line/kubectl-kustomize