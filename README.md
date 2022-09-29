
[![Go Report Card](https://goreportcard.com/badge/github.com/christianwoehrle/keycloakclient-operator)](https://goreportcard.com/report/github.com/christianwoehrle/keycloakclient-operator)

[![Coverage Status](https://coveralls.io/repos/github/christianwoehrle/keycloakclient-operator/badge.svg?branch=main)](https://coveralls.io/github/christianwoehrle/keycloakclient-operator?branch=main)

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[![codecov](https://codecov.io/gh/christianwoehrle/keycloakclient-operator/branch/main/graph/badge.svg?token=tNKcOjlxLo)](https://codecov.io/gh/christianwoehrle/keycloakclient-operator)


# KeycloakClient Operator
A Kubernetes Operator based on the Operator SDK for creating and syncing KeycloakClient-Resources in Keycloak

This Operator has it's origin from the [Legacy Keycloak Operator](https://github.com/keycloak/keycloak-operator).
If you look for the official KeycloakOperator from RedHat, please look into the [KeycloakOperator](https://github.com/keycloak/keycloak/tree/main/operator).

The Operator is opinionated in a way that it expects that Keycloak and 
the realm are already set up (i.e. with one of the available Helm Charts) and it only has 
to handle the KeycloakClients for a Keycloak Installation and a specific realm.

This fits our need as we set up Keycloak and the realm with Helm, and we have very many microservices that require their own KeycloakClient.
The Microservices are deployed via Helm and it is easy to simply deploy a KeycloakClient Resource together with the other artefacts of the Microservice and let 
the Operator handle the creation of the KeycloakClient in Keycloak.


## Try it out.

### Run the Keycloak Client Operator




To have the Keycloak Client Operator handle KeycloakClients for a specifiy Keycloak INstallation and Realm you need the following Ressources



### Keycloak CRD and Secret
You need the Keycloak-CustomResource that describes how the Keycloak Instance can be accessed (the URL) and the secret that provides Username and Password.
The Secret has to have the name of the KeycloakCRD prefixed with "credentials-"

Please see [KeycloakCR](../deploy/examples)

### Realm
The Realm-CustomResource should have id, displayName and realm set to the corresponsing name in Keycloak and the instanceSelector 
should match the labels in the KeycloakCRD.

### KeycloakClient
In the KeycloakClient you can specify the KeycloakClient. 







## Help and Documentation

* [Keycloak documentation](https://www.keycloak.org/documentation.html)
* [User Mailing List](https://groups.google.com/g/keycloak-user) - Mailing list for help and general questions about Keycloak

## Reporting an issue

If you believe you have discovered a defect in the KeycloakClent Operator please open an [an issue](https://github.com/christianwoehrle/keycloakclient-operator/issues).
Please remember to provide a good summary, description as well as steps to reproduce the issue.

## Supported Custom Resources
| *CustomResourceDefinition*                                            | *Description*                                            |
| --------------------------------------------------------------------- | -------------------------------------------------------- |
| [Keycloak](./deploy/crds/keycloak.org_keycloaks_crd.yaml)             | Manages, installs and configures Keycloak on the cluster |
| [KeycloakRealm](./deploy/crds/keycloak.org_keycloakrealms_crd.yaml)   | Represents a realm in a keycloak server                  |
| [KeycloakClient](./deploy/crds/keycloak.org_keycloakclients_crd.yaml) | Represents a client in a keycloak server                 |


## Deployment to a Kubernetes cluster


## Developer Reference
*Note*: You will need a running Kubernetes cluster to use the Operator

1. Run `make cluster/prepare` # This will apply the necessary Custom Resource Definitions (CRDs) and RBAC rules to the clusters
2. Run `kubectl apply -f deploy/operator.yaml` # This will start the operator in the current namespace

### Install keycloak with a realm 

This installs keycloak wih a realm test-realm via the codecentric helm chart

1. Run `make cluster/installKeycloak`

### Creating Keycloak Instance and realm
Once the CRDs and RBAC rules are applied and the operator is running, install the keycloak-cr, the keycloakrealm-cr and the keycloakclient-cr.
The keycloak- and keycloakrealm-crs are only used to reference keycloak and the keycloakrealm.

The keycloakclient-cr actually triggers the keycloakclient-operator to create the keycloakclient in the references keycloakcloakrealm.


1. Run `make cluster/create/examples`

<!--

### Local Development
*Note*: You will need a running Kubernetes or OpenShift cluster to use the Operator

1. clone this repo to `$GOPATH/src/github.com/keycloak/keycloak-operator`
2. run `make setup/mod cluster/prepare`
3. deploy a PostgreSQL Database -- The embedded database installation is deprecated
4. run `make code/run`
-- The above step will launch the operator on the local machine
-- To see how do debug the operator or how to deploy to a cluster, see below alternatives to step 3
5. check the IP/url of the installed Database
6. modify secret [external-db-secret.yaml](./deploy/examples/keycloak/external-db-secret.yaml) setting the values
7. execute the secret with `kubectl apply -f ./deploy/examples/keycloak/external-db-secret.yaml`
8. In a new terminal run `make cluster/create/examples`
9. Optional: configure Ingress and DNS Resolver
   - minikube: \
     -- run `minikube addons enable ingress` \
     -- run `./hack/modify_etc_hosts.sh`
   - Docker for Mac: \
     -- run `kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-0.32.0/deploy/static/provider/cloud/deploy.yaml`
        (see also https://kubernetes.github.io/ingress-nginx/deploy/) \
     -- run `./hack/modify_etc_hosts.sh keycloak.local 127.0.0.1`
10. Run `make test/e2e`

To clean the cluster (Removes CRDs, CRs, RBAC and namespace)
1. run `make cluster/clean`

#### Alternative Step 2: Debug in Goland
Debug the operator in [Goland](https://www.jetbrains.com/go/)
1. go get -u github.com/go-delve/delve/cmd/dlv
2. Create new `Go Build` debug configuration
3. Change the properties to the following
```
* Name = Keycloak Operator
* Run Kind = File
* Files = <project full path>/cmd/manager/main.go
* Working Directory = <project full path>
* Environment = KUBERNETES_CONFIG=<kube config path>;WATCH_NAMESPACE=keycloak
```
3. Apply and click Debug Keycloak operator

#### Alternative Step 3: Debug in VSCode
Debug the operator in [VS Code](https://code.visualstudio.com/docs/languages/go)
1. go get -u github.com/go-delve/delve/cmd/dlv
2. Create new launch configuration, changing your kube config location
```json
{
  "name": "Keycloak Operator",
  "type": "go",
  "request": "launch",
  "mode": "auto",
  "program": "${workspaceFolder}/cmd/manager/main.go",
  "env": {
    "WATCH_NAMESPACE": "keycloak",
    "KUBERNETES_CONFIG": "<kube config path>"
  },
  "cwd": "${workspaceFolder}",
  "args": []
}
```
3. Debug Keycloak Operator

#### Alternative Step 3: Deploying to a Cluster
Deploy the operator into the running cluster
1. build image with `operator-sdk build <image registry>/<organisation>/keycloak-operator:<tag>`. e.g. `operator-sdk build quay.io/keycloak/keycloak-operator:test`
2. Change the `image` property in `deploy/operator.yaml` to the above full image path
3. run `kubectl apply -f deploy/operator.yaml -n <NAMESPACE>`

### Makefile command reference
#### Operator Setup Management
| *Command*                      | *Description*                                                                                          |
| ------------------------------ | ------------------------------------------------------------------------------------------------------ |
| `make cluster/prepare`         | Creates the `keycloak` namespace, applies all CRDs to the cluster and sets up the RBAC files           |
| `make cluster/clean`           | Deletes the `keycloak` namespace, all `keycloak.org` CRDs and all RBAC files named `keycloak-operator` |
| `make cluster/create/examples` | Applies the example Keycloak and KeycloakRealm CRs                                                     |

#### Tests
| *Command*                    | *Description*                                               |
| ---------------------------- | ----------------------------------------------------------- |
| `make test/unit`             | Runs unit tests                                             |
| `make test/coverage/prepare` | Prepares coverage report from unit and e2e test results     |
| `make test/coverage`         | Generates coverage report                                   |

##### Running tests without cluster admin permissions
It's possible to deploy CRDs, roles, role bindings, etc. separately from running the tests:
1. Run `make cluster/prepare` as a cluster admin.
2. Run `make test/ibm-validation` as a user. The user needs the following permissions to run te tests:
```
apiGroups: ["", "apps", "keycloak.org"]
resources: ["persistentvolumeclaims", "deployments", "statefulsets", "keycloaks", "keycloakrealms", "keycloakusers", "keycloakclients", "keycloakbackups"]
verbs: ["*"]
```
Please bear in mind this is intended to be used for internal purposes as there's no guarantee it'll work without any issues.

#### Local Development
| *Command*                 | *Description*                                                                    |
| ------------------------- | -------------------------------------------------------------------------------- |
| `make setup`              | Runs `setup/mod` `setup/githooks` `code/gen`                                     |
| `make setup/githooks`     | Copys githooks from `./githooks` to `.git/hooks`                                 |
| `make setup/mod`          | Resets the main module's vendor directory to include all packages                |
| `make setup/operator-sdk` | Installs the operator-sdk                                                        |
| `make code/run`           | Runs the operator locally for development purposes                               |
| `make code/compile`       | Builds the operator                                                              |
| `make code/gen`           | Generates/Updates the operator files based on the CR status and spec definitions |
| `make code/check`         | Checks for linting errors in the code                                            |
| `make code/fix`           | Formats code using [gofmt](https://golang.org/cmd/gofmt/)                        |
| `make code/lint`          | Checks for linting errors in the code                                            |
| `make client/gen`         | Generates/Updates the clients bases on the CR status and spec definitions        |


#### CI
| *Command*           | *Description*                                                              |
| ------------------- | -------------------------------------------------------------------------- |
| `make setup/travis` | Downloads operator-sdk, makes it executable and copys to `/usr/local/bin/` |

#### Components versions

-->
## Contributing

I'm glad for any contribution. This is currently Alpha. The operator runs on my machine and I would expect that I didn't  
introduce too many errors into the orginal KeycloakOperator, as it is basically a stripped down version of the [Legacy Keycloak Operator](https://github.com/keycloak/keycloak-operator).



## Keycloak Projects

* [Keycloak](https://github.com/keycloak/keycloak) - Keycloak Server and Java adapters
* [Keycloak Documentation](https://github.com/keycloak/keycloak-documentation) - Documentation for Keycloak
* [Keycloak QuickStarts](https://github.com/keycloak/keycloak-quickstarts) - QuickStarts for getting started with Keycloak
* [Keycloak Docker](https://github.com/jboss-dockerfiles/keycloak) - Docker images for Keycloak
* [Keycloak Node.js Connect](https://github.com/keycloak/keycloak-nodejs-connect) - Node.js adapter for Keycloak
* [Keycloak Node.js Admin Client](https://github.com/keycloak/keycloak-nodejs-admin-client) - Node.js library for Keycloak Admin REST API
* [Codecentric Keycloak Helm Chart](https://artifacthub.io/packages/helm/codecentric/keycloak) - Helm chart for Keycloakx
* [Codecentric Keycloakx Helm Chart](https://artifacthub.io/packages/helm/codecentric/keycloakx) - Helm chart for Keycloakx
* [Bitnami Keycloak Helm Chart](https://github.com/bitnami/charts/tree/master/bitnami/keycloak) - Helm Chart for Keycloak
## License

* [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
