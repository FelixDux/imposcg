

# HTTP API
resource "aws_apigatewayv2_api" "imposc" {
	name          = "imposc-api"
	protocol_type = "HTTP"
	target        = aws_lambda_function.imposc.arn
}

resource "aws_apigatewayv2_integration" "imposc" {
  api_id           = aws_apigatewayv2_api.imposc.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  description               = "Serverless Impact Oscillator Model"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.imposc.invoke_arn
}

resource "aws_apigatewayv2_route" "imposc" {
  api_id    = aws_apigatewayv2_api.imposc.id
  route_key = "ANY /{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.imposc.id}"
}
