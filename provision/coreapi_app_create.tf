data "archive_file" "lambda_coreapi_apps_create" {
  type = "zip"

  source_file  = "/home/stavros/go/bin/coreapi_apps_create"
  output_path = "/home/stavros/go/bin/coreapi_apps_create.zip"
}

resource "aws_s3_object" "lambda_coreapi_apps_create" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "coreapi_apps_create.zip"
  source = data.archive_file.lambda_coreapi_apps_create.output_path

  etag = filemd5(data.archive_file.lambda_coreapi_apps_create.output_path)
}

resource "aws_lambda_function" "coreapi_apps_create" {
  function_name = "coreapi_apps_create"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_coreapi_apps_create.key

  runtime = "go1.x"
  handler = "coreapi_apps_create"

  source_code_hash = data.archive_file.lambda_coreapi_apps_create.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "coreapi_apps_create" {
  name = "/aws/lambda/${aws_lambda_function.coreapi_apps_create.function_name}"

  retention_in_days = 30
}

resource "aws_apigatewayv2_integration" "coreapi_apps_create" {
  api_id = aws_apigatewayv2_api.coreapi.id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.coreapi_apps_create.invoke_arn
}

resource "aws_apigatewayv2_route" "coreapi_apps_create" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "POST /apps"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_apps_create.id}"
}

resource "aws_lambda_permission" "apigateway_coreapi_apps_create" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.coreapi_apps_create.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.coreapi.execution_arn}/*/*/*"
}
