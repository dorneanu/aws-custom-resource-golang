import cdk = require('aws-cdk-lib');
import customResources = require('aws-cdk-lib/custom-resources');
import lambda = require('aws-cdk-lib/aws-lambda');
import { Construct } from 'constructs';

import fs = require('fs');

export interface SSMCredentialProps {
    key: string;
    value: string;
}

// SSMCredential is an AWS custom resource
//
// Example code from: https://github.com/aws-samples/aws-cdk-examples/blob/master/typescript/custom-resource/my-custom-resource.ts
export class SSMCredential extends Construct {
    public readonly response: string;

    constructor(scope: Construct, id: string, props: SSMCredentialProps) {
        super(scope, id);

        // Create lambda function from Go binary
        // TODO: Construct the lambda function either by building the golang binary
        // or uploading the binary to a S3 bucket first.
        const fn = new lambda.SingletonFunction(this, 'Singleton', {
            uuid: 'f7d4f730-4ee1-11e8-9c2d-fa7ae01bbebc',
            timeout: cdk.Duration.seconds(300),
            runtime: lambda.Runtime.GO_1_X,
        });

        // Create a new custom resource provider
        const provider = new customResources.Provider(this, 'Provider', {
            onEventHandler: fn,
        });

        const resource = new cdk.CustomResource(this, 'Resource', {
            serviceToken: provider.serviceToken,
            properties: props,
        });

        this.response = resource.getAtt('Response').toString();
    }
}
