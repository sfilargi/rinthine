resource "aws_dynamodb_table" "core_users" {
  name           = "core_users"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id_"

  attribute {
    name = "user_id_"
    type = "B"
  }
}

resource "aws_dynamodb_table" "core_users_idx_handle" {
  name           = "core_users_idx_handle"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "handle_"

  attribute {
    name = "handle_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_tokens" {
  name           = "core_tokens"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "token_"

  attribute {
    name = "token_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_tokens_idx_user_id" {
  name           = "core_tokens_idx_user_id"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id_"
  range_key      = "token_"

  attribute {
    name = "user_id_"
    type = "B"
  }

  attribute {
    name = "token_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_tokens_idx_app_id" {
  name           = "core_tokens_idx_app_id"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "app_id_"
  range_key      = "token_"

  attribute {
    name = "app_id_"
    type = "B"
  }

  attribute {
    name = "token_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_apps" {
  name           = "core_apps"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "app_id_"

  attribute {
    name = "app_id_"
    type = "B"
  }
}

resource "aws_dynamodb_table" "core_apps_idx_user_id" {
  name           = "core_apps_idx_user_id"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id_"
  range_key      = "app_id_"

  attribute {
    name = "user_id_"
    type = "B"
  }

  attribute {
    name = "app_id_"
    type = "B"
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

resource "aws_dynamodb_table" "core_oauth_codes_idx_app_id" {
  name           = "core_oauth_codes_idx_app_id"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "app_id_"
  range_key      = "code_"

  attribute {
    name = "app_id_"
    type = "B"
  }

  attribute {
    name = "code_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_oauth_codes_idx_user_id" {
  name           = "core_oauth_codes_idx_user_id"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "user_id_"
  range_key      = "code_"

  attribute {
    name = "user_id_"
    type = "B"
  }

  attribute {
    name = "code_"
    type = "S"
  }
}

resource "aws_dynamodb_table" "core_used_oauth_codes" {
  name           = "core_used_oauth_codes"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "code_"

  attribute {
    name = "code_"
    type = "S"
  }
}
