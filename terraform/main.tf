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
   s3_key    = "imposcg.zip"

   handler = "impact-oscillator"
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

resource "aws_lambda_permission" "apigw" {
   statement_id  = "AllowAPIGatewayInvoke"
   action        = "lambda:InvokeFunction"
   function_name = aws_lambda_function.imposc.function_name
   principal     = "apigateway.amazonaws.com"

   # The "/*/*" portion grants access from any method on any resource
   # within the API Gateway REST API.
   source_arn = "${aws_api_gateway_rest_api.imposc.execution_arn}/*/*"
}


# https://learn.hashicorp.com/tutorials/terraform/lambda-api-gateway#configuring-api-gateway

