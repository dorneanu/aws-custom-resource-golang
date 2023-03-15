import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as path from "path";
import { Construct } from "constructs";

// import * as sqs from 'aws-cdk-lib/aws-sqs';

export class DeploymentsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here

    // example resource
    // const queue = new sqs.Queue(this, 'DeploymentsQueue', {
    //   visibilityTimeout: cdk.Duration.seconds(300)
    // });
    //

    // Build the code and create the lambda
    //
    const environment = {
      CGO_ENABLED: "",
      GOOS: "linux",
      GOARCH: "amd64",
    };
    const lambdaPath = path.join(__dirname, "../../");

    // We could bundle the Golang binary also locally.
    // See https://github.com/aws-samples/cdk-build-bundle-deploy-example/blob/main/cdk-bundle-go-lambda-example/lib/api-stack.ts
    // But I prefer to do it in a Docker container
    const lambdaFunc = new lambda.Function(this, "GolangCustomResources", {
      code: lambda.Code.fromAsset(lambdaPath, {
        bundling: {
          image: lambda.Runtime.GO_1_X.bundlingImage,
          user: "root",
          environment,
          command: [
            "bash",
            "-c",
            [
              "cd /asset-input",
              "make build",
              "mv /asset-input/build/lambda-function.bin /asset-output/main",
            ].join(" && "),
          ],
        },
      }),
      handler: "/main",
      runtime: lambda.Runtime.GO_1_X,
    });
  }
}
