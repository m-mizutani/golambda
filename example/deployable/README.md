# Deployable Example

## Prerequisite

- Tools
  - node >= v14.7.0
  - aws-cdk >= 1.75.0
  - go >= 1.15
- Credential
  - AWS CLI credentials to deploy Lambda function via CloudFormation

## Deploy

```
$ make deploy
```

`golambda-test` will be deployed.

## Invoke lambda

Move Lambda function page in AWS console and put following JSON data as `Test Event`. Then press `Test` button.

```json
{
  "Records": [
    {
      "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
      "receiptHandle": "MessageReceiptHandle",
      "body": "{\"message\":\"blue\"}",
      "attributes": {
        "ApproximateReceiveCount": "1",
        "SentTimestamp": "1523232000000",
        "SenderId": "123456789012",
        "ApproximateFirstReceiveTimestamp": "1523232000001"
      },
      "messageAttributes": {},
      "md5OfBody": "{{{md5_of_body}}}",
      "eventSource": "aws:sqs",
      "eventSourceARN": "arn:aws:sqs:ap-northeast-1:123456789012:MyQueue",
      "awsRegion": "ap-northeast-1"
    },
    {
      "messageId": "19dd0b57-b21e-4ac1-bd88-01bbb068cb78",
      "receiptHandle": "MessageReceiptHandle",
      "body": "{\"message\":\"orange\"}",
      "attributes": {
        "ApproximateReceiveCount": "1",
        "SentTimestamp": "1523232000000",
        "SenderId": "123456789012",
        "ApproximateFirstReceiveTimestamp": "1523232000001"
      },
      "messageAttributes": {},
      "md5OfBody": "{{{md5_of_body}}}",
      "eventSource": "aws:sqs",
      "eventSourceARN": "arn:aws:sqs:ap-northeast-1:123456789012:MyQueue",
      "awsRegion": "ap-northeast-1"
    }
  ]
}
```
