# Helm (Kubernetes) plugin for drone.io

[![Build Status](http://drone.sohohousedigital.com/api/badges/ipedrazas/drone-helm/status.svg)](http://drone.sohohousedigital.com/ipedrazas/drone-helm)

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
        -e PLUGIN_HELM_COMMAND="install --name my-release stable/jenkins --debug --dry-run" \
        -e PLUGIN_API_SERVER=https://192.168.64.5:8443 \
        -e PLUGIN_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ld..." \
        -e PLUGIN_NAMESPACE=default \
        -e PLUGIN_SKIP_TLS_VERIFY=true \
        quay.io/ipedrazas/drone-helm



Happy Helming!
