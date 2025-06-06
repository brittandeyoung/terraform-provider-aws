---
subcategory: "Route 53 Domains"
layout: "aws"
page_title: "AWS: aws_route53domains_registered_domain"
description: |-
  Provides a resource to manage a domain that has been registered and associated with the current AWS account.
---

# Resource: aws_route53domains_registered_domain

Provides a resource to manage a domain that has been [registered](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/registrar-tld-list.html) and associated with the current AWS account. To register, renew and deregister a domain use the [`aws_route53domains_domain` resource](route53domains_domain.html) instead.

**This is an advanced resource** and has special caveats to be aware of when using it. Please read this document in its entirety before using this resource.

The `aws_route53domains_registered_domain` resource behaves differently from normal resources in that if a domain has been registered, Terraform does not _register_ this domain, but instead "adopts" it into management. `terraform destroy` does not delete the domain but does remove the resource from Terraform state.

## Example Usage

```terraform
resource "aws_route53domains_registered_domain" "example" {
  domain_name = "example.com"

  name_server {
    name = "ns-195.awsdns-24.com"
  }

  name_server {
    name = "ns-874.awsdns-45.net"
  }

  tags = {
    Environment = "test"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `admin_contact` - (Optional) Details about the domain administrative contact. See [Contact Blocks](#contact-blocks) for more details.
* `admin_privacy` - (Optional) Whether domain administrative contact information is concealed from WHOIS queries. Default: `true`.
* `auto_renew` - (Optional) Whether the domain registration is set to renew automatically. Default: `true`.
* `billing_contact` - (Optional) Details about the domain billing contact. See [Contact Blocks](#contact-blocks) for more details.
* `billing_privacy` - (Optional) Whether domain billing contact information is concealed from WHOIS queries. Default: `true`.
* `domain_name` - (Required) The name of the registered domain.
* `name_server` - (Optional) The list of nameservers for the domain. See [`name_server` Blocks](#name_server-blocks) for more details.
* `registrant_contact` - (Optional) Details about the domain registrant. See [Contact Blocks](#contact-blocks) for more details.
* `registrant_privacy` - (Optional) Whether domain registrant contact information is concealed from WHOIS queries. Default: `true`.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `tech_contact` - (Optional) Details about the domain technical contact. See [Contact Blocks](#contact-blocks) for more details.
* `tech_privacy` - (Optional) Whether domain technical contact information is concealed from WHOIS queries. Default: `true`.
* `transfer_lock` - (Optional) Whether the domain is locked for transfer. Default: `true`.

~> **NOTE:** You must specify the same privacy setting for `admin_privacy`, `registrant_privacy` and `tech_privacy`.

### Contact Blocks

The `admin_contact`, `billing_contact`, `registrant_contact` and `tech_contact` blocks support the following:

* `address_line_1` - (Optional) First line of the contact's address.
* `address_line_2` - (Optional) Second line of contact's address, if any.
* `city` - (Optional) The city of the contact's address.
* `contact_type` - (Optional) Indicates whether the contact is a person, company, association, or public organization. See the [AWS API documentation](https://docs.aws.amazon.com/Route53/latest/APIReference/API_domains_ContactDetail.html#Route53Domains-Type-domains_ContactDetail-ContactType) for valid values.
* `country_code` - (Optional) Code for the country of the contact's address. See the [AWS API documentation](https://docs.aws.amazon.com/Route53/latest/APIReference/API_domains_ContactDetail.html#Route53Domains-Type-domains_ContactDetail-CountryCode) for valid values.
* `email` - (Optional) Email address of the contact.
* `extra_params` - (Optional) A key-value map of parameters required by certain top-level domains.
* `fax` - (Optional) Fax number of the contact. Phone number must be specified in the format "+[country dialing code].[number including any area code]".
* `first_name` - (Optional) First name of contact.
* `last_name` - (Optional) Last name of contact.
* `organization_name` - (Optional) Name of the organization for contact types other than `PERSON`.
* `phone_number` - (Optional) The phone number of the contact. Phone number must be specified in the format "+[country dialing code].[number including any area code]".
* `state` - (Optional) The state or province of the contact's city.
* `zip_code` - (Optional) The zip or postal code of the contact's address.

### `name_server` Blocks

The `name_server` blocks supports the following:

* `glue_ips` - (Optional) Glue IP addresses of a name server. The list can contain only one IPv4 and one IPv6 address.
* `name` - (Required) The fully qualified host name of the name server.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - The domain name.
* `abuse_contact_email` - Email address to contact to report incorrect contact information for a domain, to report that the domain is being used to send spam, to report that someone is cybersquatting on a domain name, or report some other type of abuse.
* `abuse_contact_phone` - Phone number for reporting abuse.
* `creation_date` - The date when the domain was created as found in the response to a WHOIS query.
* `expiration_date` - The date when the registration for the domain is set to expire.
* `registrar_name` - Name of the registrar of the domain as identified in the registry.
* `registrar_url` - Web address of the registrar.
* `reseller` - Reseller of the domain.
* `status_list` - List of [domain name status codes](https://www.icann.org/resources/pages/epp-status-codes-2014-06-16-en).
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).
* `updated_date` - The last updated date of the domain as found in the response to a WHOIS query.
* `whois_server` - The fully qualified name of the WHOIS server that can answer the WHOIS query for the domain.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

- `create` - (Default `30m`)
- `update` - (Default `30m`)

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import domains using the domain_name. For example:

```terraform
import {
  to = aws_route53domains_registered_domain.example
  id = "example.com"
}
```

Using `terraform import`, import domains using the domain name. For example:

```console
% terraform import aws_route53domains_registered_domain.example example.com
```
