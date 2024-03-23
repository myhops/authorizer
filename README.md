# Authorizer

In our team we decided to use Istio as service mesh and make it responsible for authentication and authorization of the services.
At that time we used api keys to control access, but had planned to use OIDC and OAuth for that.

We use the Istio AuthorizationPolicy to control access.
In this resource the **when** part of the rule allows you to determine access based on the value of a given header.
Because you cannot inject the secret header value from a k8s secret, you need to put it in the policy in plain text.
This poses a security problem when you want to keep the deployment manifests in Git. 
For us this a requirement, because we use GitOps.

Istio allows you to use an external authorizer and send the relevant information over HTTP and let the authorizer decide weather access is allowed or not.
For this I created the **authorizer**.

Authorizer is a very simple web server that checks api keys contained in http headers.

## Configuration steps

To use an external authorization provider, you need to configure the Istio control plane. 
This configuration example shows the parts that you need to add to the Istio Operator manifest.

The configuration options are in the [meshConfig](https://istio.io/latest/docs/reference/config/istio.operator.v1alpha1/#IstioOperatorSpec) property of the IstioOperator.

In the [MeshConfig](https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig) you can specify the external authorizers.

This server can be used as CUSTOM validation by an Istio AuthorizationPolicy and removes the need to include secrets in the headers.

First we configure the authorizer to check on the headers **api-key=secret1** :

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    app: authorizer
  name: authorizer
spec:
  containers:
  - image: docker.io/peterzandbergen/authorizer:latest
    name: authorizer
    args:
      - key=api-key=secret1
      - logformat=json
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
---

```

So instead of checking the api key in the header:

```yaml
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: allow-on-header-value
  namespace: bar
spec:
  selector:
    matchLabels:
      app: httpbin
  action: ALLOW
  rules:
    - to:
        - operation:
      when:
        - key: request.headers[api-key]
          values:
            - valid1
            - valid2
```

you can now configure and use the authorizer:

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: ext-authz
  namespace: bar
spec:
  selector:
    matchLabels:
      app: httpbin
  action: CUSTOM
  provider:
    name: "my-custom-authz"
```

and have the authorizer configured as:

```bash
authorizer -key=api-key=valid1 -key=api-key=valid2
```

Or in a Pod

```yaml

```

```yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: example-istiocontrolplane
  namespace: istio-system
spec:
  meshConfig:
    extensionProviders:
    - envoyExtAuthzHttp:
        includeRequestHeadersInCheck:
        - X-Vbi-Api-Key
        pathPrefix: /anything
        port: 8080
        service: default/api-key-authenticator.default.svc.cluster.local
      name: api-key-authenticator
  profile: demo
```

## Istio Configuration

TODO: MeshConfiguration
TODO: Maistra operator configuration
TODO: Istio operator configuration
TODO: authorizer command




