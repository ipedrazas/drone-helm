# Helm (Kubernetes) plugin for drone.io

This plugin allows to deploy a [Helm](https://github.com/kubernetes/helm) chart into a [Kubernetes](https://github.com/kubernetes/kubernetes) cluster.

For example, this configuration will deploy Jenkins using the [stable/jenkins chart](https://github.com/kubernetes/charts/tree/master/stable/jenkins)


        pipeline:
            deploy:
                image: ipedrazas/drone-helm
                helm_command: 
                    - "install"
                    - "--name"
                    - "my-release"
                    - "stable/jenkins"
                    - "--debug" 
                    - "--dry-run"
                api_server: "http://my_api_server"
                kubernetes_token: "secret token"
                skip_tls_verify: "true"


If you don't know where to get a token from, you can execute the following command:

        kubectl exec POD_NAME -- cat /var/run/secrets/kubernetes.io/serviceaccount/token

For example, in a cluster where there's a pod called `nginx-1212390922-fdz1x` we coudl do:

        kubectl exec nginx-1212390922-fdz1x -- cat /var/run/secrets/kubernetes.io/serviceaccount/token


To test the plugin, you can run `minikube` and just run the docker image as follows:

        docker run -i ipedrazas/drone-helm /bin/drone-helm <<EOF
        {
            
            "vargs": {
                "helm_command": ["version"],
                "api_server":"https://192.168.64.5:8443",
                "kubernetes_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tcnloeTciLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjU5ZWQzYjM1LWI0MmUtMTFlNi05ZDI3LTFlZGZkMzA0MTNhNiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.aGX19xTFJdxC6d5hObKUms9Kaq0wR8fMypXTnsfjC6XjiA3_QWX9LJdRFl6wvZTRoIAjuOhAJKNAKhLQ3sK0aKrddNxn2is-HCb88BXma3BrHWOtkwajvZ1dyAhZOe2fY1I77t_mrbvTMqJ4udsom6roHf-KL8j29DJWsV0nFh6VKyvqN8f7FsNG3WuH3SFZX_LfcE0HfZxrDaVEi-CkDo0sGCqIefDk2sn4IQD6b1Ng-grJWSN-YtrcDDduEKlUHPSRMmMtWa3-Q61-yQqlyqATGbxC3UwqwaLfjCrTkg1Uikv4jWDP3-eNmuQCqG9PHKulA1riTFAgxbr09zoYxg",
                "namespace": "default",
                "kubeconfig": "/root/.kube/config",
                "skip_tls_verify": "true"
            }
        }
        EOF

This will display the helm version in the plugin and the helm version in the cluster. The following example will test 
the previous `.drone.yml` configuration: 

        docker run -i ipedrazas/drone-helm <<EOF
        {
            
            "vargs": {
                "helm_command": ["install", "--name", "my-release", "stable/jenkins", "--debug", "--dry-run"],
                "api_server":"https://192.168.64.5:8443",
                "kubernetes_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tcnloeTciLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjU5ZWQzYjM1LWI0MmUtMTFlNi05ZDI3LTFlZGZkMzA0MTNhNiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.aGX19xTFJdxC6d5hObKUms9Kaq0wR8fMypXTnsfjC6XjiA3_QWX9LJdRFl6wvZTRoIAjuOhAJKNAKhLQ3sK0aKrddNxn2is-HCb88BXma3BrHWOtkwajvZ1dyAhZOe2fY1I77t_mrbvTMqJ4udsom6roHf-KL8j29DJWsV0nFh6VKyvqN8f7FsNG3WuH3SFZX_LfcE0HfZxrDaVEi-CkDo0sGCqIefDk2sn4IQD6b1Ng-grJWSN-YtrcDDduEKlUHPSRMmMtWa3-Q61-yQqlyqATGbxC3UwqwaLfjCrTkg1Uikv4jWDP3-eNmuQCqG9PHKulA1riTFAgxbr09zoYxg",
                "namespace": "default",
                "kubeconfig": "/root/.kube/config",
                "skip_tls_verify": "true"
            }
        }
        EOF


Happy Helming!