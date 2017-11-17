# canary
Canary is a metrics and metadata exporter that runs in any cloud native environment to help you embed SLO's into your microservice application bundles.

# Building

Run the dockerfile, it will build the binary for you:

docker build -t jayunit199/canary:1.0 .

the build.sh is just a hacky convenience script for local development, use at your own risk.

# Running

Run the replication controller :)

```
kubectl create -f ./canary-rc.json
```

# Using

Go into your pod that is running, called hub-sidecar-, and run something like this:

```
kubectl exec -t -i <name of the pod> curl localhost:3000/status
```

# development policy
Canary embraces the traditional values of open source projects in the Apache and CNCF communities, and embraces ideas and community over the code itself.  Please create an issue or better yet, submit a pull request if you have any suggestions around metrics or checks that you think will be generically useful to organizations that ship code which is meant to run in a microservice environment.

# golang standards
We follow the same standards for golang as are followed in the moby project, the kubernetes project, and other major golang projects.  We embrace modern golang idioms including usage of viper (for config), glide (for dependencies), and aim to stay on the 'bleeding edge', since, after all, we aim to always deploy inside of containers.
