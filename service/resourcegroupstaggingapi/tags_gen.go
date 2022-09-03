// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package resourcegroupstaggingapi

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	tftags "github.com/infracasts/terraform-provider-aws-public/tags"
)

// []*SERVICE.Tag handling

// Tags returns resourcegroupstaggingapi service tags.
func Tags(tags tftags.KeyValueTags) []*resourcegroupstaggingapi.Tag {
	result := make([]*resourcegroupstaggingapi.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &resourcegroupstaggingapi.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from resourcegroupstaggingapi service tags.
func KeyValueTags(tags []*resourcegroupstaggingapi.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}
