## cnvrg.io slim operator delivery chart
Run `make` or follow the steps below to manually build the chart dependencies.

1. Clone down the lastest slim operator. We need the cnvrg-operator chart local until the slim operator
   is pushed to the Chart Mueseum.
   `git clone -b DEV-15047-slim-cnvrg  git@github.com:AccessibleAI/cnvrg-operator.git`
   
2. You need to change the name of the cnvrg operator chart in cnvrg-operator/charts/operator/Chart.yaml to
   cnvrg-operator. The reason is both the app and operator chart have the same name which won't work. Here
   is how the operator/Chart.yaml should look.
   ```
   apiVersion: v2
   name: cnvrg-operator
   description: A cnvrg.io operator v3 chart for K8s
   type: application
   version: 4.3.22-DEV-15047-slim-cnvrg-101
   appVersion: 1.2.3
   ```
   
4. Change into the slim-operator-delivery directory, so we can run the helm install command and update the packages.
   ```
   cd slim-operator-delivery
   ```
   
5. Update the build dependency charts. This downloads all of the charts locally for installatation. You can view all the 
   local charts downloaded in the ./charts directory.
   ```
   helm dependency update
   ```
   
6. Update the values.yaml. You can add any values for the charts under the approperiate key in the values file. For example, 
   to add an externalIP address to the ingress nginx deployment, the yaml file would look as follows:
   ```
   ingress-nginx: 
     # -- Set to true to enabled deployment of nginx ingress 
     enabled: true
     controller:
      service:
        externalIPs:
          - 172.31.21.239
   ```
   **Note:** You can always get the chart values by running the following command against the chart located in `./charts`
`helm show values ./charts/<chart-file-name>`

7. Set the registry password which you will pass later to the helm install command.
  `export PASSWORD=<cnvrghelm-dockerhub-password>`

8. Helm install command:
   ```
   helm upgrade cnvrg -n cnvrg . --wait -f values.yaml --install --create-namespace \
   --set cnvrg-operator.registry.password=$PASSWORD \
   --set cnvrg.registry.password=$PASSWORD
   ```

### Ingress Controller installs
You have the ability to use this chart to deploy the ingress controllers separately.
Here are some examples:

#### Ingress nginx
```
helm upgrade nginx -n nginx . --wait --create-namespace --install \
--set ingress-nginx.enabled=true \
--set cnvrg-operator.enabled=false \
--set ingress-nginx.controller.service.externalIPs={172.31.21.239}
```
#### Istio gateway
```
helm upgrade ingressgateway -n istio-system . --wait --create-namespace --install \
--set base.enabled=true \
--set cnvrg-operator.enabled=false \
--set gateway.service.externalIPs={172.31.21.239}
```

