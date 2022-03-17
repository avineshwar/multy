package object_storage

import (
	"github.com/multycloud/multy/resources/common"
)

// aws_s3_bucket
type AwsS3Bucket struct {
	*common.AwsResource `hcl:",squash" default:"name=aws_s3_bucket"`
	Bucket              string `hcl:"bucket"`
}
