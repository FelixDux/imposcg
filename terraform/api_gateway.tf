# resource "aws_api_gateway_rest_api" "imposc" {
#   name        = "ImpactOscillator"
#   description = "Serverless Impact Oscillator Model"
#   binary_media_types = ["image/png"]
# }

# # Route all incoming requests
# resource "aws_api_gateway_resource" "proxy" {
#    rest_api_id = aws_api_gateway_rest_api.imposc.id
#    parent_id   = aws_api_gateway_rest_api.imposc.root_resource_id
#    path_part   = "{proxy+}"
# }

# resource "aws_api_gateway_method" "proxy" {
#    rest_api_id   = aws_api_gateway_rest_api.imposc.id
#    resource_id   = aws_api_gateway_resource.proxy.id
#    http_method   = "ANY"
#    authorization = "NONE"
#    api_key_required = false
# }

# # route to the lambda function
# resource "aws_api_gateway_integration" "lambda" {
#    rest_api_id = aws_api_gateway_rest_api.imposc.id
#    resource_id = aws_api_gateway_method.proxy.resource_id
#    http_method = aws_api_gateway_method.proxy.http_method

#    integration_http_method = "POST"
#    type                    = "AWS_PROXY"
#    uri                     = aws_lambda_function.imposc.invoke_arn
# }

# resource "random_id" "id" {
# 	byte_length = 8
# }

# HTTP API
resource "aws_apigatewayv2_api" "imposc" {
	name          = "imposc-api" #${random_id.id.hex}"
	protocol_type = "HTTP"
	target        = aws_lambda_function.imposc.arn
	# route_key     = "ANY /{proxy+}" 
	# route_key     = "$default"
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

# # The proxy resource cannot match an empty path at the root of the API, so apply a configuration at the root resource
# resource "aws_api_gateway_method" "proxy_root" {
#    rest_api_id   = aws_api_gateway_rest_api.imposc.id
#    resource_id   = aws_api_gateway_rest_api.imposc.root_resource_id
#    http_method   = "ANY"
#    authorization = "NONE"
# }

# resource "aws_api_gateway_integration" "lambda_root" {
#    rest_api_id = aws_api_gateway_rest_api.imposc.id
#    resource_id = aws_api_gateway_method.proxy_root.resource_id
#    http_method = aws_api_gateway_method.proxy_root.http_method

#    integration_http_method = "POST"
#    type                    = "AWS_PROXY"
#    uri                     = aws_lambda_function.imposc.invoke_arn
# }

# # Activate the configuration and expose the API at a URL that can be used for testing:
# resource "aws_api_gateway_deployment" "impact_oscillator" {
#    depends_on = [
#      aws_api_gateway_integration.lambda,
#      aws_api_gateway_integration.lambda_root,
#    ]

#    rest_api_id = aws_api_gateway_rest_api.imposc.id
#    stage_name  = "default"
# }

# output "base_url" {
#   value = aws_api_gateway_deployment.impact_oscillator.invoke_url
# }


