resource "aws_dynamodb_table" "users" {
  name           = "users"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "handle"

  attribute {
    name = "handle"
    type = "S"
  }
}