import "source-map-support/register";
import * as cdk from '@aws-cdk/core';
import * as lambda from '@aws-cdk/aws-lambda';

import * as path from 'path';

export class GolambdaTestStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // --------------------------------------
    // Lambda
    new lambda.Function(this, 'handler', {
      runtime: lambda.Runtime.GO_1_X,
      handler: 'handler',
      code: lambda.Code.fromAsset(path.join(__dirname, '../build')),
      timeout: cdk.Duration.seconds(10),
      memorySize: 128,
      reservedConcurrentExecutions: 1,
      environment: {
        SENTRY_DSN: process.env.SENTRY_DSN || "",
      },
    });
  }
}

const app = new cdk.App();
new GolambdaTestStack(app, "golambda-test");
