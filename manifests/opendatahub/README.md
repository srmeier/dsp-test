# Data Science Pipelines

Data Science Pipelines is the Open Data Hub's pipeline solution for data scientists. It is built on top of the upstream [Kubeflow Piplines](https://github.com/kubeflow/pipelines) and [kfp-tekton](https://github.com/kubeflow/kfp-tekton) projects. The Open Data Hub community has a [fork](https://github.com/opendatahub-io/data-science-pipelines) of this upstream under the Open Data Hub org.


## Installation

### Prerequisites

1. The cluster needs to be OpenShift 4.9 or higher
2. OpenShift Pipelines 1.7.2 or higher needs to be installed on the cluster
3. The Open Data Hub operator needs to be installed
4. The default installation namespace for Data Science Pipelines is `odh-applications`. This namespace will need to be created. In case you wish to install in a custom location, create it and update the kfdef as documented below.

### Installation Steps

1. Ensure that the prerequisites are met.
2. Apply the kfdef at [kfctl_openshift_ds-pipelines.yaml](https://github.com/opendatahub-io/odh-manifests/blob/master/kfdef/kfctl_openshift_ds-pipelines.yaml). You may need to update the `namespace` field under `metadata` in case you want to deploy in a namespace that isn't `odh-applications`.
3. To find the url for Data Science pipelines, you can run the following command.
    ```bash
    $ oc get route -n <kdef_namespace> ds-pipeline-ui -o jsonpath='{.spec.host}'
    ```
    The value of `<kfdef_namespace>` should match the namespace field of the kfdef that you applied.
4. Alternatively, you can access the route via the console. To do so:

    1. Go to `<kfdef_namespace>`
    2. Click on `Networking` in the sidebar on the left side.
    3. Click on `Routes`. It will take you to a new page in the console.
    4. Click the url under the `Location` column for the row item matching `ds-pipeline-ui`


## Directory Structure

### Base

This directory contains artifacts for deploying all backend components of Data Science Pipelines. This deployment currently includes the kfp-tekton backend as well as a Minio deployment to act as an object store. The Minio deployment will be moved to an overlay at some point in the near future.

### Overlays

1. metadata-store-mariadb: This overlay contains artifacts for deploying a MariaDB database. MySQL-based databases are currently the only supported backend for Data Science Pipelines, so if you don't have an existing MySQL database deployed, this overlay can be applied to satisfy the requirement.
2. metadata-store-mysql: This overlay contains artifacts for deploying a MySQL database. MySQL-based databases are currently the only supported backend for Data Science Pipelines, so if you don't have an existing MySQL database deployed, this overlay can be applied to satisfy the requirement.
3. metadata-store-postgresql: This overlay contains artifacts for deploying a PostgreSQL database. Data Science Pipelines does not currently support PostgreSQL as a backend, so deploying this overlay will not actually modify Data Science Pipelines behaviour.
4. ds-pipeline-ui: This overlay contains deployment artifacts for the Data Science Pipelines UI. Deploying Data Science Pipelines without this overlay will result in only the backend artifacts being created.
5. object-store-minio: This overlay contains artifacts for deploying Minio as the Object Store to store Pipelines artifacts.
6. default-configs: This overlay creates ConfigMaps and Secrets with default values for a deployment with both a local MySQL database and Minio object store. *Note*: Using this overlay allows for a simple and quick setup, but also marks the configs as managed objects when used with the ODH Operator, which will reconcile any post-deployment changes made, and cannot be overridden.
7. integration-odhdashboard: Adds resources required to integrate the Data Science Pipelines application into the ODH Dashboard UI, such as documentation and application launcher tiles.
8. component-mlmd: Adds the ML-Metadata component which provides artifact lineage tracking in the UI.

### Prometheus

This directory contains the service monitor definition for Data Science Pipelines. It is always deployed by base, so this will eventually be moved into the base directory itself.

## Parameters

You can customize the Data Science Pipelines deployment by injecting custom parameters to change the default deployment. The following parameters can be used:

* **pipeline_install_configuration**: The ConfigMap name that contains the values to install the Data Science Pipelines environment. This parameter defaults to `pipeline-install-config` and you can find an example in the [repository](./base/configmaps/pipeline-install-config.yaml).
* **ds_pipelines_configuration**: The ConfigMap name that contains the values to integrate Data Science Pipelines with the underlying components (Database and Object Store). This parameter defaults to `kfp-tekton-config` and you can find an example in the [repository](./base/configmaps/kfp-tekton-config.yaml).
* **database_secret**: The secret that contains the credentials for the Data Science Pipelines Databse. It defaults to `mysql-secret` if using the `metadata-store-mysql` overlay or `postgresql-secret` if using the `metadata-store-postgresql` overlay.
* **ds_pipelines_ui_configuration**: The ConfigMap that contains the values to customize UI. It defaults to `ds-pipeline-ui-configmap`.

## Configuration

* It is possible to configure what S3 storage is being used by Pipeline Runs. Detailed instructions on how to configure this will be added once Minio is moved to an overlay.

## Usage

### These instructions will be updated once Data Science Pipelines has a tile available in odh-dashboard

1. Go to the ds-pipelines-ui route.
2. Click on `Pipelines` on the left side.
3. There will be a `[Demo] flip-coin` Pipeline already available. Click on it.
4. Click on the blue `Create run` button towards the top of the screen.
5. You can leave all the fields untouched. If desired, you can create a new experiment to link the pipeline run to, or rename the run itself.
6. Click on the blue `Start` button.
7. You will be taken to the `Runs` page. You will see a row matching the `Run name` you previously picked. Click on the `Run name` in that row.
8. Once the Pipeline is done running, you can see a graph of all the pods that were created as well as the paths that were followed.
9. For further verification, you can view all the pods that were created as part of the Pipeline Run in the `<kfdef_namespace>`. They will all show up as `Completed`.

## Data Science Pipelines Architecture

A complete architecture can be found at [ODH Data Science Pipelines Architecture and Design](https://docs.google.com/document/d/1o-JS1uZKLZsMY3D16kl5KBdyBb-aV-kyD_XycdJOYpM/edit#heading=h.3aocw3evrps0). This document will be moved to GitHub once the corresponding ML Ops SIG repos are created.
