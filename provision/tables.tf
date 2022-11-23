resource "aws_dynamodb_table" "core_userhandles" {
  name           = "core_user_handles"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "handle"

  attribute {
    name = "handle"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_users" {
  name           = "core_users"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "userid"

  attribute {
    name = "userid"
    type = "B"
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

resource "aws_dynamodb_table" "core_app_names" {
  name           = "core_app_names"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "name"

  attribute {
    name = "name"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_apps" {
  name           = "core_apps"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "appid"

  attribute {
    name = "appid"
    type = "B"
  }
}

resource "aws_dynamodb_table" "core_user_apps" {
  name           = "core_user_apps"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "userid"
  range_key      = "appid"

  attribute {
    name = "userid"
    type = "B"
  }

  attribute {
    name = "appid"
    type = "B"
  }
}


