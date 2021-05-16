provider "aws" {
    profile = "default"
    region = "eu-west-2"
}

resource "aws_s3_bucket" "staging" {
    bucket = "felixdux-imposcg-lambda"
    acl = "private"
    force_destroy = "true"
}

resource "aws_lambda_function" "imposc" {
   function_name = "ImpactOscillator"

   s3_bucket = "felixdux-imposcg-lambda"
   s3_key    = "v1.0.0/imposcg.zip"

   # "main" is the filename within the zip file (main.js) and "handler"
   # is the name of the property under which the handler function was
   # exported in that file.
   handler = "impact-oscillator"
   runtime = "go1.16.x"

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

# https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway#configuring-api-gateway

