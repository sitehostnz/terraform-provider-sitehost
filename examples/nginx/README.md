# SiteHost Linux VPS + Nginx

This example launches a Ubuntu 22.04 LTS server with a nginx web server.

> Note: To run this example, first configure your SiteHost provider as described in <https://kb.sitehost.nz/developers/terraform-provider>.

## Prerequisites

Your `Client ID` and `API Key` values can be generated<https://cp.sitehost.nz/api/list-keys> by clicking "Add API Key". You can find out how to do this in our [knowledge base](https://kb.sitehost.nz/developers/api).

You will need to export your SiteHost Client ID and API Key as an environment variable:

```sh
export TF_VAR_sitehost_client_id="Put Your SiteHost Client ID Here"
export TF_VAR_sitehost_api_key="Put Your SiteHost API key Here"
```

## Run this example

From the `examples/nginx` directory.

```sh
terraform init
terraform apply
```

> The server installation should be completed in under 5 minutes.

## Destroy the Resources

Clean up by removing all the resources that were created in one command:

```sh
terraform destroy
```
