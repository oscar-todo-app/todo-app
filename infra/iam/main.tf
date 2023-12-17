terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}


resource "aws_iam_policy_document" "todo-secrets" {
  statement {
    actions = [
      "secretsmanager:GetSecretValue",
      "secretsmanager:DescribeSecret"
    ]
    effect = "Allow"
    resources = [
      var.secret_arn
    ]
  }

}
resource "aws_iam_policy" "todo-secrets" {
  name   = "todo-secrets"
  policy = data.aws_iam_policy_document.todo-secrets.json

}

resource "aws_iam_role" "todo-secrets" {
  name = "todo-secrets"
  assume_role_policy = jsondecode({

    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "VisualEditor0",
        "Effect" : "Allow",
        "Action" : [
          "secretsmanager:GetSecretValue",
          "secretsmanager:DescribeSecret"
        ],
        "Resource" : var.secretarn
      }
    ]
  })

}

resource "aws_iam_openid_connect_provider" "todo-secrets" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = var.thumbprint_list
  url             = var.iamurl

}
resource "aws_iam_" "name" {

}
