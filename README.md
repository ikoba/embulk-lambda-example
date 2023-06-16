# Emublk Lambda Example

This is an example of running [Embulk](https://www.embulk.org/) on AWS Lambda.

## How to deploy

1. Edit `src/main.go` to suit your purpose.

2. Build a docker image.

    `docker build -t embulk .`

    If you build on Apple Silicon, Gem installation processes are probably unstable.
    In that case, update Docker Desktop to the latest version, open Settings menu, enable the option below.

    - Settings > Features in development

        ✅ Use Rosetta for x86/amd64 emulation on Apple Silicon

3. Push the docker image to an Amazon ECR Repository.

    Please see below for more information.

    [Pushing a Docker image](https://docs.aws.amazon.com/AmazonECR/latest/userguide/docker-push-ecr-image.html)

4. Create a Lambda function.

5. Configure the Lambda function.

    - Deploy the image pushed to Amazon ECR in the above step.

    - Increase the memory size allocated to the Lambda function to 512MB or more.
      If it is executed with the default memory size (128MB), it will probably fail.

    - Set the timeout long enough to execute the function.

    - Set [permissions](#Permissions) properly.

    - Select a VPC, Subnets and Security Groups properly.

## Permissions

The Execution Role must be set properly in order for the Lambda function to be executed successfully.

At least the following policy statement is required for log output (This is set by default).

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "logs:CreateLogGroup",
            "Resource": "arn:aws:logs:<your-region>:<your-account>:*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": [
                "arn:aws:logs:<your-region>:<your-account>:log-group:/aws/lambda/<your-lambda-function>:*"
            ]
        }
    ]
}
```

To access resources such as RDS or Redshift in a VPC, the following statement is also required.

```json
        {
            "Effect": "Allow",
            "Action": [
                "ec2:CreateNetworkInterface",
                "ec2:DescribeNetworkInterfaces",
                "ec2:DeleteNetworkInterface"
            ],
            "Resource": "*"
        }

```

To use embulk-output-redshift plugin like this example, the following statement is also required.

```json
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::<your-backet>",
                "arn:aws:s3:::<your-backet>/*"
            ]
        }
```

To retrieve parameters like DB password from AWS SSM Parameter Store, the following statement is also required.

```json
        {
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameter",
                "kms:Decrypt"
            ],
            "Resource": "*"
        }
```

## Notes

### Which value should be set to `auth_method` / `aws_auth_method`?

When using a plugin that requires authentication to AWS such as follows, the value of `auth_method` / `aws_auth_method` in a config file should be `env`.

- embulk-input-s3
- embulk-output-s3
- embulk-output-redshift

If `env` is set, the plugin refers environment variables such as AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, and AWS_SESSION_TOKEN.

This means that access permissions associated with the Lambda function's execution role will be applied.

Please see below for more information.

- [Working with Lambda execution environment credentials](https://docs.aws.amazon.com/lambda/latest/dg/lambda-intro-execution-role.html#permissions-executionrole-source-function-arn)

- [S3 file input plugin for Embulk](https://github.com/embulk/embulk-input-s3#configuration)

### Embulk v0.11

The Embulk v0.11 has various changes from the previous stable version v0.9.

Please see below for the chnages.

[Embulk v0.11 is coming soon](https://www.embulk.org/articles/2023/04/13/embulk-v0.11-is-coming-soon.html1)
