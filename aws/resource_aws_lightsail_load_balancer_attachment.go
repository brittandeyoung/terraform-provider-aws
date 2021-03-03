package aws

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/lightsail/waiter"
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
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
		},
	}
}

func resourceAwsLightsailLoadBalancerAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn
	instanceNames := expandInstanceNames(d.Get("instance_names").(*schema.Set))

	req := lightsail.AttachInstancesToLoadBalancerInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
		InstanceNames:    instanceNames,
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

	_, err = waiter.OperationCreated(conn, op.Id)
	if err != nil {
		return fmt.Errorf("Error waiting for load balancer attatchment (%s) to become ready: %s", d.Id(), err)
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
				log.Printf("[WARN] Lightsail load balancer (%s) not found, removing Attachment from state", d.Id())
				d.SetId("")
				return nil
			}
			return err
		}
		return err
	}

	lbhs := resp.LoadBalancer.InstanceHealthSummary

	var iNames []string

	if len(lbhs) > 0 {
		for _, element := range lbhs {
			iNames = append(iNames, aws.StringValue(element.InstanceName))
		}
	} else {
		log.Printf("[WARN] Lightsail load balancer (%s) has no attachments, removing from state", d.Id())
		d.SetId("")
	}

	d.Set("instance_names", iNames)

	return nil
}

func resourceAwsLightsailLoadBalancerAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn
	instanceNames := expandInstanceNames(d.Get("instance_names").(*schema.Set))

	req := lightsail.DetachInstancesFromLoadBalancerInput{
		LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
		InstanceNames:    instanceNames,
	}

	resp, err := conn.DetachInstancesFromLoadBalancer(&req)

	if err != nil {
		return err
	}

	if len(resp.Operations) == 0 {
		return fmt.Errorf("No operations found for DetachInstancesFromLoadBalancer request")
	}

	op := resp.Operations[0]

	_, err = waiter.OperationCreated(conn, op.Id)
	if err != nil {
		return fmt.Errorf("Error waiting for load balancer attatchment (%s) to become detached: %s", d.Id(), err)
	}

	return err
}

func resourceAwsLightsailLoadBalancerAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lightsailconn

	if d.HasChange("instance_names") {
		instanceNames := expandInstanceNames(d.Get("instance_names").(*schema.Set))
		req := lightsail.AttachInstancesToLoadBalancerInput{
			LoadBalancerName: aws.String(d.Get("load_balancer_name").(string)),
			InstanceNames:    instanceNames,
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

		_, err = waiter.OperationCreated(conn, op.Id)
		if err != nil {
			return fmt.Errorf("Error waiting for load balancer attatchment (%s) to become ready: %s", d.Id(), err)
		}

		if err != nil {
			return err
		}
	}

	return resourceAwsLightsailLoadBalancerAttachmentRead(d, meta)
}

func expandInstanceNames(d *schema.Set) []*string {
	// instanceNamesSet := d.Get("instance_names").(*schema.Set)
	instanceNames := make([]*string, d.Len())
	for i, instanceName := range d.List() {
		instanceNames[i] = aws.String(instanceName.(string))
	}
	return instanceNames
}
