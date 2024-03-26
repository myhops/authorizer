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

This README explains how to use it with the Istio Operator and in the Istio Authorization Policy.

## Configuration and installation steps

We need to install the authorizer and configure Istio and add an envoyExtAuthzHttp to the extensionProviders.
Secondly we need to install the authorizer service with the set of valid API keys.
And lastly we create an AuthorizationPolicy with a custom action that uses the external authorizer.

### Istio Operator Configuration

Below you see the manifest for the Istio operator and the part that adds the extension.

```yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: example-istiocontrolplane
  namespace: istio-system
spec:
  meshConfig:
    extensionProviders:     # 1
      - envoyExtAuthzHttp:  # 2
          name: api-key-authenticator   # 3
          includeRequestHeadersInCheck: # 4
            - X-Vbi-Api-Key
          # pathPrefix: / # Add a prefix if necessary
          port: 8080
          # The DNS name of the service.
          service: api-key-authenticator.default.svc.cluster.local # 5
  profile: demo
```

1. The mesh config allows you to add multiple extensions.
1. We add the http type authorizer.
1. The name of the extension. We need this in the AP.
1. The header(s) that contain the api key.
1. The authenticator service DNS

### The AuthorizationPolicy

Now that we have configured Istio with the HTTP Authorizer extension, we can use it in the Authorization Policy resource.

Assume that you have a workload **httpbin** that exposes an HTTP service with path /anything that you are only allowed to access with a valid API key.
The key needs to be in the header **X-Api-Key** but we do not want to expose the value of the key. 

With the extension we can configure this as follows:

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: ext-authz
spec:
  selector:
    matchLabels:
      app: httpbin # 1
  action: CUSTOM   # 2
  provider:
    name: api-key-authenticator # 3
```

1. The selector of the workload we want to protect.
1. This indicates that we use a custom authorizer.
1. The name we gave the extension.

### Authorizer

We want to configure the authorizer with the keys in a Secret.
This enables us to use e.g. SealedSecrets and store these in the same repo as the kustomization files for our application. 
The Authorizer accepts it configuration via cli flags, but we will show how we can use environment variables for this.

The secret contains the following info.
Please note that we use stringData instead of data for readability.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: auth-keys
stringData:
  keys: X-Vbi-Api-Key=key1,X-Vbi-Api-Key=key2
```

We can now write our manifest for the Authorizer service. 
We inject the api keys as env variables.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: authorizer
  name: authorizer
spec:
  replicas: 1
  selector:
    matchLabels:
      app: authorizer
  template:
    metadata:
      labels:
        app: authorizer
    spec:
      containers:
      - image: docker.io/peterzandbergen/authorizer:v0.0.2
        imagePullPolicy: IfNotPresent
        name: authorizer
        envFrom:
          - secretRef:
              name: auth-keys # 1
            prefix: AUTH_     # 2
        args:
          - -key=$(AUTH_keys) # 3
          - -logformat=json
        ports:
          - name: http
            containerPort: 8080
            protocol: TCP  
        resources: 
          limits:
            cpu: 100m
            memory: 200Mi
          requests:
            cpu: 50m
            memory: 100Mi
```

1. Get the auth-keys in the environment.
1. Use a prefix to prevent collisions.
1. Set the -key flag.

The Authorizer can accept multiple keys in one -key flag or multiple -key flags.

Now we have everything we need in place to use API keys to protect our service without having to adapt our workload.













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


you can now configure and use the authorizer:


and have the authorizer configured as:

```bash
authorizer -key=api-key=valid1 -key=api-key=valid2
```

Or in a Pod

```yaml

```



## Istio Configuration

- TODO: MeshConfiguration
- TODO: Maistra operator configuration
- TODO: Istio operator configuration
- TODO: authorizer command




