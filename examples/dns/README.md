# SiteHost DNS

This example adds a [DNS Zone](https://kb.sitehost.nz/dns/dns-zones) to SiteHost Control Panel and adds the MX records to link your domain with [Google Workspace](https://support.google.com/a/answer/174125?hl=en).

> Note: To run this example, first configure your SiteHost provider as described in <https://kb.sitehost.nz/developers/terraform-provider>.

## Prerequisites

Your `Client ID` and `API Key` values can be generated<https://cp.sitehost.nz/api/list-keys> by clicking "Add API Key". You can find out how to do this in our [knowledge base](https://kb.sitehost.nz/developers/api).

You will need to export your SiteHost Client ID and API Key as an environment variable:

```sh
export TF_VAR_sitehost_client_id="Put Your SiteHost Client ID Here"
export TF_VAR_sitehost_api_key="Put Your SiteHost API key Here"
```

## Set your Domain Name

Edit the `domain_name` value in the [variable.tf](./variable.tf) file.

## Run this example

From the `examples/dns` directory.

```sh
terraform init
terraform apply
```

> To use the DNS Zones manager you need to use the [SiteHost Nameservers](https://kb.sitehost.nz/domains/name-servers).

## Destroy the Resources

Clean up by removing all the resources that were created in one command:

```sh
terraform destroy
```
