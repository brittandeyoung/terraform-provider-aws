package aws

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsLightsailLoadBalancerAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsLightsailLoadBalancerAttachmentCreate,
		Read:   resourceAwsLightsailLoadBalancerAttachmentRead,
		Update: resourceAwsLightsailLoadBalancerAttachmentUpdate,
		Delete: resourceAwsLightsailLoadBalancerAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(2, 255),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z]`), "must begin with an alphabetic character"),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_\-.]+[^._\-]$`), "must contain only alphanumeric characters, underscores, hyphens, and dots"),
				),
			},
			"instance_names": {
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
		},
	}
}

func resourceAwsLightsailLoadBalancerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn

	req := lightsail.AttachInstancesToLoadBalancerInput{
		InstanceNames:  aws.String(d.Get("instance_names").(string)),
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
	}

	resp, err := conn.AttachInstancesToLoadBalancer(&req)
	if err != nil {
		return err
	}

	if len(resp.Operations) == 0 {
		return fmt.Errorf("No operations found for AttachInstancesToLoadBalancer request")
	}

	op := resp.Operations[0]
	d.SetId(d.Get("load_balancer_name").(string))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Started"},
		Target:     []string{"Completed", "Succeeded"},
		Refresh:    resourceAwsLightsailLoadBalancerAttachmentOperationRefreshFunc(op.Id, meta),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		// We don't return an error here because the Create call succeeded
		log.Printf("[ERR] Error waiting for load balancer attachment (%s) to become ready: %s", d.Id(), err)
	}

	return resourceAwsLightsailLoadBalancerAttachmentRead(d, meta)
}

func resourceAwsLightsailLoadBalancerAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn

	resp, err := conn.GetLoadBalancer(&lightsail.GetLoadBalancerInput{
		LoadBalancerName: aws.String(d.Id()),
	})

	
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NotFoundException" {
				log.Printf("[WARN] Lightsail load balancer Attachment (%s) not found, removing from state", d.Id())
				d.SetId("")
				return nil
			}
			return err
		}
		return err
	}
  for index, element := range someSlice {

	var s []string
	for index, element := range resp.LoadBalancer.InstanceHealthSummary {
		s = append(s, )
	
	d.Set("instance_names", resp.LoadBalancer.Arn)
	d.Set("created_at", resp.LoadBalancer.CreatedAt.Format(time.RFC3339))
	d.Set("health_check_path", resp.LoadBalancer.HealthCheckPath)
	d.Set("instance_port", resp.LoadBalancer.InstancePort)
	d.Set("name", resp.LoadBalancer.Name)
	d.Set("protocol", resp.LoadBalancer.Protocol)
	d.Set("public_ports", resp.LoadBalancer.PublicPorts)
	d.Set("dns_name", resp.LoadBalancer.DnsName)

	if err := d.Set("tags", keyvaluetags.LightsailKeyValueTags(resp.LoadBalancer.Tags).IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsLightsailLoadBalancerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn
	resp, err := conn.DeleteLoadBalancer(&lightsail.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(d.Id()),
	})

	op := resp.Operations[0]

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Started"},
		Target:     []string{"Completed", "Succeeded"},
		Refresh:    resourceAwsLightsailLoadBalancerAttachmentOperationRefreshFunc(op.Id, meta),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for load balancer (%s) to become destroyed: %s",
			d.Id(), err)
	}

	return err
}

func resourceAwsLightsailLoadBalancerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn

	if d.HasChange("health_check_path") {
		_, err := conn.UpdateLoadBalancerAttribute(&lightsail.UpdateLoadBalancerAttributeInput{
			AttributeName:    aws.String("HealthCheckPath"),
			AttributeValue:   aws.String(d.Get("health_check_path").(string)),
			LoadBalancerName: aws.String(d.Get("name").(string)),
		})
		d.Set("health_check_path", d.Get("health_check_path").(string))
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")

		if err := keyvaluetags.LightsailUpdateTags(conn, d.Id(), o, n); err != nil {
			return fmt.Errorf("error updating Lightsail Instance (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceAwsLightsailLoadBalancerAttachmentRead(d, meta)
}

// method to check the status of an Operation, which is returned from
// Create/Delete methods.
// Status's are an aws.OperationStatus enum:
// - NotStarted
// - Started
// - Failed
// - Completed
// - Succeeded (not documented?)
func resourceAwsLightsailLoadBalancerAttachmentOperationRefreshFunc(
	oid *string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn := meta.(*AWSClient).lightsailconn
		log.Printf("[DEBUG] Checking if Lightsail Operation (%s) is Completed", *oid)
		o, err := conn.GetOperation(&lightsail.GetOperationInput{
			OperationId: oid,
		})
		if err != nil {
			return o, "FAILED", err
		}

		if o.Operation == nil {
			return nil, "Failed", fmt.Errorf("Error retrieving Operation info for operation (%s)", *oid)
		}

		log.Printf("[DEBUG] Lightsail Operation (%s) is currently %q", *oid, *o.Operation.Status)
		return o, *o.Operation.Status, nil
	}
}
