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


data "aws_iam_policy_document" "todo-secrets" {
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



module "iam_secret-role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"

  create_role                   = true
  role_name                     = "todo-secret"
  provider_url                  = replace(var.provider_url, "https://", "")
  role_policy_arns              = [aws_iam_policy.todo-secrets.arn]
  oidc_fully_qualified_subjects = ["system:serviceaccount:todo:secret-sa"]
}
