package amp

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/infracasts/terraform-provider-aws-public/conns"
	tftags "github.com/infracasts/terraform-provider-aws-public/tags"
)

func DataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prometheus_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).AMPConn
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	workspaceID := d.Get("workspace_id").(string)
	workspace, err := FindWorkspaceByID(conn, workspaceID)

	if err != nil {
		return fmt.Errorf("reading AMP Workspace (%s): %w", workspaceID, err)
	}

	d.SetId(workspaceID)

	d.Set("alias", workspace.Alias)
	d.Set("arn", workspace.Arn)
	d.Set("created_date", workspace.CreatedAt.Format(time.RFC3339))
	d.Set("prometheus_endpoint", workspace.PrometheusEndpoint)
	d.Set("status", workspace.Status.StatusCode)

	if err := d.Set("tags", KeyValueTags(workspace.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	return nil
}
