package opensearch

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/opensearchservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/infracasts/terraform-provider-aws-public/conns"
	"github.com/infracasts/terraform-provider-aws-public/tfresource"
	"github.com/infracasts/terraform-provider-aws-public/verify"
)

func ResourceDomainPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainPolicyUpsert,
		Read:   resourceDomainPolicyRead,
		Update: resourceDomainPolicyUpsert,
		Delete: resourceDomainPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(180 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_policies": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: verify.SuppressEquivalentPolicyDiffs,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
		},
	}
}

func resourceDomainPolicyRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).OpenSearchConn

	ds, err := FindDomainByName(conn, d.Get("domain_name").(string))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] OpenSearch Domain Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading OpenSearch Domain Policy (%s): %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Received OpenSearch domain: %s", ds)

	policies, err := verify.PolicyToSet(d.Get("access_policies").(string), aws.StringValue(ds.AccessPolicies))

	if err != nil {
		return err
	}

	d.Set("access_policies", policies)

	return nil
}

func resourceDomainPolicyUpsert(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).OpenSearchConn
	domainName := d.Get("domain_name").(string)

	policy, err := structure.NormalizeJsonString(d.Get("access_policies").(string))

	if err != nil {
		return fmt.Errorf("policy (%s) is invalid JSON: %w", policy, err)
	}

	_, err = conn.UpdateDomainConfig(&opensearchservice.UpdateDomainConfigInput{
		DomainName:     aws.String(domainName),
		AccessPolicies: aws.String(policy),
	})
	if err != nil {
		return err
	}

	d.SetId("esd-policy-" + domainName)

	if err := waitForDomainUpdate(conn, d.Get("domain_name").(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for OpenSearch Domain Policy (%s) to be updated: %w", d.Id(), err)
	}

	return resourceDomainPolicyRead(d, meta)
}

func resourceDomainPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).OpenSearchConn

	_, err := conn.UpdateDomainConfig(&opensearchservice.UpdateDomainConfigInput{
		DomainName:     aws.String(d.Get("domain_name").(string)),
		AccessPolicies: aws.String(""),
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for OpenSearch domain policy %q to be deleted", d.Get("domain_name").(string))

	if err := waitForDomainUpdate(conn, d.Get("domain_name").(string), d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for OpenSearch Domain Policy (%s) to be deleted: %w", d.Id(), err)
	}

	return nil
}