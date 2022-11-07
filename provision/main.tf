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

resource "aws_cloudwatch_log_group" "coreapi_gateway" {
  name = "/aws/api_gw/${aws_api_gateway_rest_api.coreapi.name}"

  retention_in_days = 30
}

resource "aws_api_gateway_rest_api" "coreapi" {
  name        = "CoreAPI"
  description = "Rinthine Core API"
}

resource "aws_api_gateway_resource" "coreapi_users" {
  rest_api_id = aws_api_gateway_rest_api.coreapi.id
  parent_id   = aws_api_gateway_rest_api.coreapi.root_resource_id
  path_part   = "users"
}

resource "aws_api_gateway_method" "coreapi_users_create" {
  rest_api_id   = aws_api_gateway_rest_api.coreapi.id
  resource_id   = aws_api_gateway_resource.coreapi_users.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "coreapi_users_create" {
  rest_api_id = aws_api_gateway_rest_api.coreapi.id
  resource_id = aws_api_gateway_method.coreapi_users_create.resource_id
  http_method = aws_api_gateway_method.coreapi_users_create.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.coreapi_users_create.invoke_arn
}

resource "aws_api_gateway_deployment" "coreapi_production" {
  depends_on = [
    aws_api_gateway_integration.coreapi_users_create,
  ]

  rest_api_id = "${aws_api_gateway_rest_api.coreapi.id}"
  stage_name  = "production"
}

resource "aws_api_gateway_base_path_mapping" "api" {
  depends_on = [
    aws_api_gateway_deployment.coreapi_production
  ]
  api_id      = aws_api_gateway_rest_api.coreapi.id
  stage_name  = aws_api_gateway_deployment.coreapi_production.stage_name
  domain_name = "api.rinthine.com"
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name =  aws_lambda_function.coreapi_users_create.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/* portion grants access from any method on any resource
  # within the API Gateway "REST API".
  source_arn = "${aws_api_gateway_rest_api.coreapi.execution_arn}/*/*/*"
}

