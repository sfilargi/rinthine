terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.2.0"
    }
  }

  required_version = "~> 1.0"
}

provider "aws" {
  region = var.aws_region
}

resource "random_id" "lambda_bucket_name" {
  byte_length = 16
}

resource "aws_s3_bucket" "lambda_bucket" {
  bucket = "lambda-${random_id.lambda_bucket_name.hex}"
}

resource "aws_s3_bucket_acl" "bucket_acl" {
  bucket = aws_s3_bucket.lambda_bucket.id
  acl    = "private"
}

data "archive_file" "lambda_coreapi_users_create" {
  type = "zip"

  source_file  = "/home/stavros/go/bin/coreapi_users_create"
  output_path = "/home/stavros/go/bin/coreapi_users_create.zip"
}

resource "aws_s3_object" "lambda_coreapi_users_create" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "coreapi_users_create.zip"
  source = data.archive_file.lambda_coreapi_users_create.output_path

  etag = filemd5(data.archive_file.lambda_coreapi_users_create.output_path)
}

resource "aws_lambda_function" "coreapi_users_create" {
  function_name = "coreapi_users_create"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_coreapi_users_create.key

  runtime = "go1.x"
  handler = "coreapi_users_create"

  source_code_hash = data.archive_file.lambda_coreapi_users_create.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "coreapi_users_create" {
  name = "/aws/lambda/${aws_lambda_function.coreapi_users_create.function_name}"

  retention_in_days = 30
}

data "archive_file" "lambda_coreapi_login" {
  type = "zip"

  source_file  = "/home/stavros/go/bin/coreapi_login"
  output_path = "/home/stavros/go/bin/coreapi_login.zip"
}

resource "aws_s3_object" "lambda_coreapi_login" {
  bucket = aws_s3_bucket.lambda_bucket.id

  key    = "coreapi_login.zip"
  source = data.archive_file.lambda_coreapi_login.output_path

  etag = filemd5(data.archive_file.lambda_coreapi_login.output_path)
}

resource "aws_lambda_function" "coreapi_login" {
  function_name = "coreapi_login"

  s3_bucket = aws_s3_bucket.lambda_bucket.id
  s3_key    = aws_s3_object.lambda_coreapi_login.key

  runtime = "go1.x"
  handler = "coreapi_login"

  source_code_hash = data.archive_file.lambda_coreapi_login.output_base64sha256

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_cloudwatch_log_group" "coreapi_login" {
  name = "/aws/lambda/${aws_lambda_function.coreapi_login.function_name}"

  retention_in_days = 30
}


resource "aws_iam_role" "lambda_exec" {
  name = "serverless_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
          Service = "lambda.amazonaws.com"
      	}
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "dynamodb_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}


resource "aws_apigatewayv2_api" "coreapi" {
  name          = "coreapi"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "coreapi_production" {
  api_id = aws_apigatewayv2_api.coreapi.id

  name        = "production"
  auto_deploy = true

  access_log_settings {
     destination_arn = aws_cloudwatch_log_group.apigateway_coreapi.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      caller                  = "$context.identity.caller"
      userARN                 = "$context.identity.userArn"
      }
    )
  }
}

resource "aws_apigatewayv2_integration" "coreapi_users_create" {
  api_id = aws_apigatewayv2_api.coreapi.id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.coreapi_users_create.invoke_arn
}

resource "aws_apigatewayv2_route" "coreapi_users_create" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "POST /users"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_users_create.id}"
}

resource "aws_apigatewayv2_integration" "coreapi_login" {
  api_id = aws_apigatewayv2_api.coreapi.id

  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.coreapi_login.invoke_arn
}

resource "aws_apigatewayv2_route" "coreapi_login" {
  api_id = aws_apigatewayv2_api.coreapi.id

  route_key = "POST /login"
  target    = "integrations/${aws_apigatewayv2_integration.coreapi_login.id}"
}

resource "aws_apigatewayv2_domain_name" "core_rinthine_com" {
  domain_name = "core.rinthine.com"

  domain_name_configuration {
    certificate_arn = "arn:aws:acm:us-east-1:562555332644:certificate/6be46e8e-6278-4c33-b4db-b314e4f7e8c7"
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

data "aws_route53_zone" "rinthine_com" {
  name = "rinthine.com."
  private_zone = false
}

resource "aws_route53_record" "core" {
  name    = aws_apigatewayv2_domain_name.core_rinthine_com.domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.rinthine_com.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.core_rinthine_com.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.core_rinthine_com.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_api_mapping" "core_rinthine_com" {
  api_id      = aws_apigatewayv2_api.coreapi.id
  domain_name = aws_apigatewayv2_domain_name.core_rinthine_com.id
  stage       = aws_apigatewayv2_stage.coreapi_production.id
}

resource "aws_cloudwatch_log_group" "apigateway_coreapi" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.coreapi.name}"

  retention_in_days = 30
}

resource "aws_lambda_permission" "apigateway_coreapi_users_create" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.coreapi_users_create.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.coreapi.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "apigateway_coreapi_login" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.coreapi_login.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.coreapi.execution_arn}/*/*/*"
}
