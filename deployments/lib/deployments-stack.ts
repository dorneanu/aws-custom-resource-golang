import * as cdk from "aws-cdk-lib";
import * as path from "path";
import * as customResources from "aws-cdk-lib/custom-resources";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { spawnSync, SpawnSyncOptions } from "child_process";
import { Construct } from "constructs";
import { SSMCredential } from "./custom-resource";

export class DeploymentsStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Build the Golang based Lambda function
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
          // try to bundle on the local machine
          // From https://github.com/aws-samples/cdk-build-bundle-deploy-example/tree/
          local: {
            tryBundle(outputDir: string) {
              // make sure that we have all the required
              // dependencies to build the executable locally.
              // In this case we just check to make sure we have
              // go installed
              try {
                exec("go version", {
                  stdio: [
                    // show output
                    "ignore", //ignore stdio
                    process.stderr, // redirect stdout to stderr
                    "inherit", // inherit stderr
                  ],
                });
              } catch {
                // if we don't have go installed return false which
                // tells the CDK to try Docker bundling
                return false;
              }

              exec(
                [
                  "make build", // run tests first
                  `GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${path.join(
                    outputDir,
                    "main"
                  )} ./cmd/main.go`,
                  // `go build -mod=vendor -o ${path.join(
                  //   outputDir,
                  //   "bootstrap"
                  // )}`,
                ].join(" && "),
                {
                  env: { ...process.env, ...environment }, // environment variables to use when running the build command
                  stdio: [
                    // show output
                    "ignore", //ignore stdio
                    process.stderr, // redirect stdout to stderr
                    "inherit", // inherit stderr
                  ],
                  cwd: lambdaPath, // where to run the build command from
                }
              );
              return true;
            },
          },
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

    // Create a new custom resource provider
    const provider = new customResources.Provider(this, "Provider", {
      onEventHandler: lambdaFunc,
    });

    // Create custom resource
    new SSMCredential(this, "SSMCredential1", provider, {
      key: "/testing",
      value: "some-secret-value",
    });

    new SSMCredential(this, "SSMCredential2", provider, {
      key: "/testing2",
      value: "some-secret-value",
    });

    new SSMCredential(this, "SSMCredential3", provider, {
      key: "/testing3",
      value: "some-secret-value",
    });
    //
    new SSMCredential(this, "SSMCredential4", provider, {
      key: "/testing4",
      value: "some-secret-value",
    });

    new SSMCredential(this, "SSMCredential5", provider, {
      key: "/testing5",
      value: "some-secret-value",
    });

    new SSMCredential(this, "SSMCredential6", provider, {
      key: "/testing6",
      value: "some-secret-value",
    });

    new SSMCredential(this, "SSMCredential7", provider, {
      key: "/testing6",
      value: "some-secret-value",
    });
  }
}

// Utility function
function exec(command: string, options?: SpawnSyncOptions) {
  const proc = spawnSync("bash", ["-c", command], options);

  if (proc.error) {
    throw proc.error;
  }

  if (proc.status != 0) {
    if (proc.stdout || proc.stderr) {
      throw new Error(
        `[Status ${proc.status}] stdout: ${proc.stdout
          ?.toString()
          .trim()}\n\n\nstderr: ${proc.stderr?.toString().trim()}`
      );
    }
    throw new Error(`go exited with status ${proc.status}`);
  }

  return proc;
}
