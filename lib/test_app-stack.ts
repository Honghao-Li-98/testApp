import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as path from "path";

// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class TestAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create an API Gateway REST API
    const api = new apigateway.RestApi(this, "MyRestApi", {
      restApiName: "My API",
      description: "My API Description",
      deployOptions: {
        stageName: "prod",
      },
    });

    // GOOS=linux GOARCH=amd64 go build -o main main.go

    // Create a Lambda function
    const lambdaFunction = new lambda.Function(this, "MyLambdaFunction", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../src/functions/main.zip")
      ),
      handler: "main",
    });

    // Add a resource and a GET method to the API
    const resource = api.root.addResource("hello");
    const method = resource.addMethod(
      "GET",
      new apigateway.LambdaIntegration(lambdaFunction)
    );
  }
}
