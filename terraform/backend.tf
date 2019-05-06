terraform {
  backend "s3" {
    bucket = "papi-tfstate-backend"
    key    = "tfstate/terraform.tfstate"
    region = "eu-west-2"
    dynamodb_table = "papi-tfstate-lock"
  }
}