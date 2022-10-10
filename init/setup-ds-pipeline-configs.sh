app_namespace=${1:-odh-applications}

config_params=""
# Check defaultable envvars to print info message
defaultable_envvars="S3_ENDPOINT S3_ACCESS_KEY S3_SECRET_KEY DB_HOST DB_PORT DB_USERNAME DB_PASSWORD DB_DATABASE S3_BUCKET"
for envvar in $defaultable_envvars; do
    if [[ -z "${!envvar}" ]]; then
        echo "Environment Variable '${envvar}' not set, using default value."
    else
        config_params="${config_params} -p ${envvar}=${!envvar}"
    fi
done

oc process -f ../manifests/ds-pipelines-config-templates.yaml ${config_params} | oc apply  -n ${app_namespace} -f -
