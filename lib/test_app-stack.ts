import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as path from "path";
import { Secret } from "aws-cdk-lib/aws-secretsmanager";
import { Effect, PolicyStatement } from "aws-cdk-lib/aws-iam";

export class TestAppStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const secret = Secret.fromSecretNameV2(this, "test-secret", "test-secret");

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
    const saveFn = new lambda.Function(this, "MyLambdaFunction", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../src/functions/save/main.zip")
      ),
      handler: "main",
      environment: {
        SECRET_ARN: secret.secretName,
      },
      initialPolicy: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: [
            "secretsmanager:GetResourcePolicy",
            "secretsmanager:GetSecretValue",
            "secretsmanager:DescribeSecret",
            "secretsmanager:ListSecretVersionIds",
            "secretsmanager:ListSecrets",
          ],
          resources: ["*"],
        }),
      ],
    });

    const getFn = new lambda.Function(this, "myget", {
      runtime: lambda.Runtime.GO_1_X,
      code: lambda.Code.fromAsset(
        path.join(__dirname, "../src/functions/get/main.zip")
      ),
      handler: "main",
      environment: {
        SECRET_ARN: secret.secretName,
      },
      initialPolicy: [
        new PolicyStatement({
          effect: Effect.ALLOW,
          actions: [
            "secretsmanager:GetResourcePolicy",
            "secretsmanager:GetSecretValue",
            "secretsmanager:DescribeSecret",
            "secretsmanager:ListSecretVersionIds",
            "secretsmanager:ListSecrets",
          ],
          resources: ["*"],
        }),
      ],
    });

    // Add a resource and a GET method to the API
    const resource = api.root.addResource("hello");
    resource.addMethod("POST", new apigateway.LambdaIntegration(saveFn));
    resource.addMethod("GET", new apigateway.LambdaIntegration(getFn));
  }
}
