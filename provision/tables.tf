resource "aws_dynamodb_table" "core_users" {
  name           = "core_users"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "handle_"

  attribute {
    name = "handle_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_user_tokens" {
  name           = "core_user_tokens"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "token_"

  attribute {
    name = "token_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_apps" {
  name           = "core_apps"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "name_"

  attribute {
    name = "name_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_user_apps" {
  name           = "core_user_apps"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_"
  range_key      = "app_"

  attribute {
    name = "user_"
    type = "S"
  }

  attribute {
    name = "app_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_oauth_codes" {
  name           = "core_oauth_codes"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "code_"

  attribute {
    name = "code_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_oauth_used_codes" {
  name           = "core_oauth_used_codes"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "code_"

  attribute {
    name = "code_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_app_codes" {
  name           = "core_app_codes"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "app_"
  range_key      = "code_"

  attribute {
    name = "app_"
    type = "S"
  }

  attribute {
    name = "code_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_app_tokens" {
  name           = "core_app_tokens"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "app_"
  range_key      = "token_"

  attribute {
    name = "app_"
    type = "S"
  }

  attribute {
    name = "token_"
    type = "S"
  }
}
