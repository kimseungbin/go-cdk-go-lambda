package main

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"log"
	"os"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ImageResizerStackProps struct {
	awscdk.StackProps
}

func NewImageResizerStack(scope constructs.Construct, id string, props *ImageResizerStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &stackProps)

	// The code that defines your stack goes here

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("ImageResizerQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

	return stack
}

var requiredEnvVars = []string{"ORIGINAL_BUCKET_ARN"}

func main() {
	err := checkRequiredEnvVars(requiredEnvVars)
	if err != nil {
		log.Fatal(err)
	}

	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewImageResizerStack(app, "ImageResizerStack", &ImageResizerStackProps{
		awscdk.StackProps{
			Env: env(),
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
