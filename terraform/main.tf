provider "aws" {
    profile = "default"
    region = "eu-west-2"
}

resource "aws_s3_bucket" "staging" {
    bucket = "felixdux-imposcg-lambda"
    acl = "private"
    force_destroy = "true"
}

variable source_archive {
  type        = string
  default     = "imposcg.zip"
  description = "S3 key for the zip file containing the source code"
}

resource "aws_lambda_function" "imposc" {
   function_name = "ImpactOscillator"

   s3_bucket = "felixdux-imposcg-lambda"
   s3_key    = var.source_archive

   handler = "impact.oscillator"
   runtime = "go1.x"

   role = aws_iam_role.lambda_exec.arn
}

 # IAM role which dictates what other AWS services the Lambda function
 # may access.
resource "aws_iam_role" "lambda_exec" {
   name = "impact_oscillator_lambda"

   assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

}

# Add AWSLambdaBasicExecutionRole managed policy to role so that the function can write logs to CloudWatch

resource "aws_iam_policy" "lambda_logs" {
  name        = "lambda_logs"
  description = "Write to CloudWatch"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:PutLogEvents"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_attach" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_logs.arn
}

# Grant API gateway permission to invoke

resource "aws_lambda_permission" "apigw" {
   statement_id  = "AllowAPIGatewayInvoke"
   action        = "lambda:InvokeFunction"
   function_name = aws_lambda_function.imposc.function_name
   principal     = "apigateway.amazonaws.com"

   # The "/*/*" portion grants access from any method on any resource
   # within the API Gateway REST API.
   source_arn = "${aws_apigatewayv2_api.imposc.execution_arn}/*/*"
}


# https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway#configuring-api-gateway

