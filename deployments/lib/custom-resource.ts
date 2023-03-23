import * as path from "path";
import * as cdk from "aws-cdk-lib";
import * as customResources from "aws-cdk-lib/custom-resources";
import { Construct } from "constructs";
import fs = require("fs");

export interface SSMCredentialProps {
  key: string;
  value: string;
}

// SSMCredential is an AWS custom resource
//
// Example code from: https://github.com/aws-samples/aws-cdk-examples/blob/master/typescript/custom-resource/my-custom-resource.ts
export class SSMCredential extends Construct {
  public readonly response: string;

  constructor(
    scope: Construct,
    id: string,
    provider: customResources.Provider,
    props: SSMCredentialProps
  ) {
    super(scope, id);

    const resource = new cdk.CustomResource(this, id, {
      serviceToken: provider.serviceToken,
      properties: props,
    });

    this.response = resource.getAtt("Response").toString();
  }
}
