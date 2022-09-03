// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/aws/aws-sdk-go/service/directoryservice/directoryserviceiface"
	tftags "github.com/infracasts/terraform-provider-aws-public/tags"
)

// ListTags lists ds service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn directoryserviceiface.DirectoryServiceAPI, identifier string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier)
}

func ListTagsWithContext(ctx context.Context, conn directoryserviceiface.DirectoryServiceAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &directoryservice.ListTagsForResourceInput{
		ResourceId: aws.String(identifier),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns ds service tags.
func Tags(tags tftags.KeyValueTags) []*directoryservice.Tag {
	result := make([]*directoryservice.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &directoryservice.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from directoryservice service tags.
func KeyValueTags(tags []*directoryservice.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates ds service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn directoryserviceiface.DirectoryServiceAPI, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn directoryserviceiface.DirectoryServiceAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &directoryservice.RemoveTagsFromResourceInput{
			ResourceId: aws.String(identifier),
			TagKeys:    aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.RemoveTagsFromResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &directoryservice.AddTagsToResourceInput{
			ResourceId: aws.String(identifier),
			Tags:       Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.AddTagsToResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
