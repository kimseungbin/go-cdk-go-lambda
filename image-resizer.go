package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"log"
	"os"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoCdkStackProps struct {
	awscdk.StackProps
}

func NewGoCdkStack(scope constructs.Construct, id string, props *GoCdkStackProps) (stack awscdk.Stack) {
	var stackProps awscdk.StackProps

	if props != nil {
		stackProps = props.StackProps
	}

	stack = awscdk.NewStack(scope, &id, &stackProps)

	// Get the original bucket by ARN from the environment variable
	originalBucketArn := os.Getenv("ORIGINAL_BUCKET_ARN")
	awss3.Bucket_FromBucketArn(stack, jsii.String("OriginalBucket"), jsii.String(originalBucketArn))

	// Create a role for the Lambda function
	lambdaRole := awsiam.NewRole(stack, jsii.String("LambdaRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})

	// Create a policy for the role allowing read access to the original bucket
	originalBucketReadAccessPolicy := awsiam.NewPolicy(stack, jsii.String("OriginalBucketReadAcessPolicy"), &awsiam.PolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					jsii.String("s3:GetObject"),
				},
				Effect: awsiam.Effect_ALLOW,
				Resources: &[]*string{
					jsii.String(originalBucketArn + "/*"),
				},
			}),
		},
	})

	lambdaRole.AttachInlinePolicy(originalBucketReadAccessPolicy)

	awscdklambdagoalpha.NewGoFunction(stack, jsii.String("ImageProcessor"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   jsii.String("lambda"),
	})

	return
}

var requiredEnvVars = []string{"ORIGINAL_BUCKET_ARN"}

func main() {
	err := checkRequiredEnvVars(requiredEnvVars)
	if err != nil {
		log.Fatal(err)
	}

	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoCdkStack(app, "ImageProcessorStack", &GoCdkStackProps{
		awscdk.StackProps{
			Env: nil,
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}

// checkRequiredEnvVars checks that all required environment variables are set.
func checkRequiredEnvVars(requiredEnvVars []string) (err error) {
	var missingVars []string

	for _, envVar := range requiredEnvVars {
		if _, isExist := os.LookupEnv(envVar); !isExist {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missingVars)
	}

	return
}
