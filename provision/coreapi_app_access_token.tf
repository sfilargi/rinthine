data "archive_file" "lambda_coreapi_access_token" {
  type = "zip"

  source_file  = "/home/stavros/go/bin/coreapi_access_token"
  output_path = "/home/stavros/go/bin/coreapi_access_token.zip"
}

resource "aws_s3_object" "lambda_coreapi_access_token" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "coreapi_access_token.zip"
  source = data.archive_file.lambda_coreapi_access_token.output_path

  etag = filemd5(data.archive_file.lambda_coreapi_access_token.output_path)
}

resource "aws_lambda_function" "coreapi_access_token" {
  function_name = "coreapi_access_token"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_coreapi_access_token.key

  runtime = "go1.x"
  handler = "coreapi_access_token"

  source_code_hash = data.archive_file.lambda_coreapi_access_token.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "coreapi_access_token" {
  name = "/aws/lambda/${aws_lambda_function.coreapi_access_token.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "coreapi_access_token" {
  api_id = aws_apigatewayv2_api.coreapi.id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.coreapi_access_token.invoke_arn
}

resource "aws_apigatewayv2_route" "coreapi_access_token_get" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "GET /access_token"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_access_token.id}"
}

resource "aws_apigatewayv2_route" "coreapi_access_token_post" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "POST /access_token"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_access_token.id}"
}

resource "aws_lambda_permission" "apigateway_coreapi_access_token" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.coreapi_access_token.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.coreapi.execution_arn}/*/*/*"
}
