# Helm (Kubernetes) plugin for drone.io

[![Build Status](https://drone.sohohousedigital.com/api/badges/ipedrazas/drone-helm/status.svg)](https://drone.sohohousedigital.com/ipedrazas/drone-helm)

This plugin allows to deploy a [Helm](https://github.com/kubernetes/helm) chart into a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster.

For example, this configuration will deploy my-app using the [stable/jenkins chart](https://github.com/kubernetes/charts/tree/master/stable/jenkins)


        pipeline:
             helm_deploy:
                image: quay.io/ipedrazas/drone-helm                    
                chart: stable/jenkins
                release: my-dear-jenkins

There are two secrets you have to create:

                drone secret add --image=quay.io/ipedrazas/drone-helm \
                        your-user/your-repo API_SERVER https://mykubernetesapiserver


                drone secret add --image=quay.io/ipedrazas/drone-helm \
                        your-user/your-repo KUBERNETES_TOKEN eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJ...

                        
If you don't know where to get a token from, you can execute the following command:

        kubectl exec POD_NAME -- cat /var/run/secrets/kubernetes.io/serviceaccount/token

For example, in a cluster where there's a pod called `nginx-1212390922-fdz1x` we coudl do:

        kubectl exec nginx-1212390922-fdz1x -- cat /var/run/secrets/kubernetes.io/serviceaccount/token


To test the plugin, you can run `minikube` and just run the docker image as follows:


        docker run --rm \
        -e PLUGIN_API_SERVER=https://192.168.64.5:8443 \
        -e PLUGIN_TOKEN="" \
        -e PLUGIN_NAMESPACE=default \
        -e PLUGIN_SKIP_TLS_VERIFY=true \
        -e PLUGIN_RELEASE=my-release \
        -e PLUGIMN_CHART=stable/redis \
        -e PLUGIN_VALUES="tag=TAG,api=API" \
        -e PLUGIN_SECRETS=TAG,API \
        -e PLUGIN_DEBUG=true \
        -e PLUGIN_DRY_RUN=true \
        -e DRONE_BUILD_EVENT=delete \
        quay.io/ipedrazas/drone-helm


## Secrets

If you find that you need to put a secret in the `--set` values of your `Helm` command you have to create the drone secret first:

                drone secret add --image=quay.io/ipedrazas/drone-helm \
                        your-user/your-repo MYSECRET secretvalue

Then you have to define values as 


                pipeline:
                  helm_deploy:
                    image: quay.io/ipedrazas/drone-helm                    
                    chart: stable/jenkins
                    release: my-dear-jenkins
                    values: webhook.token=${MYSECRET},webhook.key=$KEY
                    api_server: ${STAGING_API_SERVER}
                    secrets: MYSECRET,STAGING_API_SERVER,KEY

You have to do this because from 0.5 version fo Drone, secrets are not expanded in plugins. This means that there's
no possibility of passing secret parameters as part of a value to the plugin.

This is a limitation of Drone. to overcome that problem, we define the `SECRETS` and the plugin will resolve them

Happy Helming!
