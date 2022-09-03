// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package memorydb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/aws/aws-sdk-go/service/memorydb/memorydbiface"
	tftags "github.com/infracasts/terraform-provider-aws-public/tags"
)

// ListTags lists memorydb service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn memorydbiface.MemoryDBAPI, identifier string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier)
}

func ListTagsWithContext(ctx context.Context, conn memorydbiface.MemoryDBAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &memorydb.ListTagsInput{
		ResourceArn: aws.String(identifier),
	}

	output, err := conn.ListTagsWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.TagList), nil
}

// []*SERVICE.Tag handling

// Tags returns memorydb service tags.
func Tags(tags tftags.KeyValueTags) []*memorydb.Tag {
	result := make([]*memorydb.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &memorydb.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from memorydb service tags.
func KeyValueTags(tags []*memorydb.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates memorydb service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn memorydbiface.MemoryDBAPI, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn memorydbiface.MemoryDBAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &memorydb.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			TagKeys:     aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &memorydb.TagResourceInput{
			ResourceArn: aws.String(identifier),
			Tags:        Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
