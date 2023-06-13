# Image Resizer

This CDK deploys stacks that resize images to FHD resolution. It requires two ARNs of S3 buckets. When a file is added to the original bucket, an event notification is triggered, which then invokes the resizer Lambda function. Subsequently, the resizer Lambda function uploads the resized image to the corresponding directory in the resized bucket.



# How to deploy

- set environment variables `ORIGINAL_BUCKET_ARN` and `RESIZED_BUCKET_ARN`.
- cdk deploy

# How it works

Loads buckets using ARNs
```Go
originalBucketArn := os.Getenv("ORIGINAL_BUCKET_ARN")
originalBucket := awss3.Bucket_FromBucketArn(stack, jsii.String("OriginalBucket"), jsii.String(originalBucketArn))

resizedBucketArn := os.Getenv("RESIZED_BUCKET_ARN")
resizedBucket := awss3.Bucket_FromBucketArn(stack, jsii.String("ResizedBucket"), jsii.String(resizedBucketArn))
```

Creates Lambda function
```Go
lambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("ImageProcessor"), &awscdklambdagoalpha.GoFunctionProps{
    Runtime: awslambda.Runtime_GO_1_X(),
    Entry:   jsii.String("lambda"),
    Environment: &map[string]*string{
        "RESIZED_BUCKET_NAME": resizedBucket.BucketName(),
    },
    Timeout: awscdk.Duration_Seconds(jsii.Number(900)),
})
```

Add trigger to the Lambda function
```Go
lambdaDestination := awss3notifications.NewLambdaDestination(lambda)
originalBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, lambdaDestination)
```

Give permissions to the Lambda function
```Go
originalBucket.GrantRead(lambda, "*")
resizedBucket.GrantReadWrite(lambda, "*")
```

# Things to refactor

- handle failing logs separately from success logs
- notify failing logs in slack