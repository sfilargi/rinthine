data "archive_file" "lambda_coreapi_authorize" {
  type = "zip"

  source_file  = "/home/stavros/go/bin/coreapi_authorize"
  output_path = "/home/stavros/go/bin/coreapi_authorize.zip"
}

resource "aws_s3_object" "lambda_coreapi_authorize" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "coreapi_authorize.zip"
  source = data.archive_file.lambda_coreapi_authorize.output_path

  etag = filemd5(data.archive_file.lambda_coreapi_authorize.output_path)
}

resource "aws_lambda_function" "coreapi_authorize" {
  function_name = "coreapi_authorize"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_coreapi_authorize.key

  runtime = "go1.x"
  handler = "coreapi_authorize"

  source_code_hash = data.archive_file.lambda_coreapi_authorize.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "coreapi_authorize" {
  name = "/aws/lambda/${aws_lambda_function.coreapi_authorize.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "coreapi_authorize" {
  api_id = aws_apigatewayv2_api.coreapi.id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.coreapi_authorize.invoke_arn
}

resource "aws_apigatewayv2_route" "coreapi_authorize_get" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "GET /authorize"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_authorize.id}"
}

resource "aws_apigatewayv2_route" "coreapi_authorize_post" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "POST /authorize"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_authorize.id}"
}

resource "aws_lambda_permission" "apigateway_coreapi_authorize" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.coreapi_authorize.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.coreapi.execution_arn}/*/*/*"
}
