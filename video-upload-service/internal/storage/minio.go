package storage

import (
    "context"
    "fmt"
    "mime/multipart"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type MinioClient struct {
    client     *s3.Client
    bucketName string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string) (*MinioClient, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion("us-east-1"),
        config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(context.Context) (aws.Credentials, error) {
            return aws.Credentials{
                AccessKeyID:     accessKey,
                SecretAccessKey: secretKey,
                Source:          "minio",
            }, nil
        })),
        config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
            return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
        })),
    )
    if err != nil {
        return nil, err
    }

    return &MinioClient{
        client:     s3.NewFromConfig(cfg),
        bucketName: bucket,
    }, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, file multipart.File, fileName string) (string, error) {
    _, err := m.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(m.bucketName),
        Key:    aws.String(fileName),
        Body:   file,
    })
    if err != nil {
        return "", err
    }
    return fmt.Sprintf("%s/%s", m.bucketName, fileName), nil
}
