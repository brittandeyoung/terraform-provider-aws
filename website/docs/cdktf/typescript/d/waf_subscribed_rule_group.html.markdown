---
subcategory: "WAF Classic"
layout: "aws"
page_title: "AWS: aws_waf_subscribed_rule_group"
description: |-
  Retrieves information about a Managed WAF Rule Group from AWS Marketplace.
---


<!-- Please do not edit this file, it is generated. -->
# Data Source: aws_waf_subscribed_rule_group

`aws_waf_subscribed_rule_group` retrieves information about a Managed WAF Rule Group from AWS Marketplace (needs to be subscribed to first).

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { Token, TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { DataAwsWafSubscribedRuleGroup } from "./.gen/providers/aws/data-aws-waf-subscribed-rule-group";
import { WafWebAcl } from "./.gen/providers/aws/waf-web-acl";
interface MyConfig {
  defaultAction: any;
  metricName: any;
  name: any;
}
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string, config: MyConfig) {
    super(scope, name);
    const byMetricName = new DataAwsWafSubscribedRuleGroup(
      this,
      "by_metric_name",
      {
        metricName: "F5BotDetectionSignatures",
      }
    );
    const byName = new DataAwsWafSubscribedRuleGroup(this, "by_name", {
      name: "F5 Bot Detection Signatures For AWS WAF",
    });
    new WafWebAcl(this, "acl", {
      rules: [
        {
          priority: 1,
          ruleId: Token.asString(byName.id),
          type: "GROUP",
        },
        {
          priority: 2,
          ruleId: Token.asString(byMetricName.id),
          type: "GROUP",
        },
      ],
      defaultAction: config.defaultAction,
      metricName: config.metricName,
      name: config.name,
    });
  }
}

```

## Argument Reference

This data source supports the following arguments:

* `name` - (Optional) Name of the WAF rule group.
* `metricName` - (Optional) Name of the WAF rule group.

At least one of `name` or `metricName` must be configured.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `id` - ID of the WAF rule group.

<!-- cache-key: cdktf-0.20.8 input-a3e30b303725ac3ef31fa0ab75097ff240470ba2637fa75b3c06ac4caf9ba1ba -->