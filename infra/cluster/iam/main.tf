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

resource "aws_iam_policy" "cert-manager" {
  name        = "cert-manager"
  description = "cert-manager policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : "route53:GetChange",
        "Resource" : "arn:aws:route53:::change/*"
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "route53:ChangeResourceRecordSets",
          "route53:ListResourceRecordSets"
        ],
        "Resource" : "arn:aws:route53:::hostedzone/*"
      },
      {
        "Effect" : "Allow",
        "Action" : "route53:ListHostedZonesByName",
        "Resource" : "*"
      },
    ]
  })
}

module "iam_secret-role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-assumable-role-with-oidc"

  create_role                   = true
  role_name                     = "todo-secret"
  provider_url                  = replace(var.provider_url, "https://", "")
  role_policy_arns              = [aws_iam_policy.todo-secrets.arn]
  oidc_fully_qualified_subjects = ["system:serviceaccount:todo:secret-sa"]
}


module "cert_manager_irsa_role" {
  source                     = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  create_role                = true
  attach_cert_manager_policy = true
  role_name                  = "todo-cert-manager"
  oidc_providers = {
    eks = {
      provider_arn               = var.provider_arn
      namespace_service_accounts = ["cert-manager:cert-manager"]
    }
  }


}

module "external_dns_irsa_role" {
  source = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"

  role_name                     = "external-dns"
  attach_external_dns_policy    = true
  external_dns_hosted_zone_arns = ["arn:aws:route53:::hostedzone/Z05080821D3KFPK0X4CL1"]

  oidc_providers = {
    eks = {
      provider_arn               = var.provider_arn
      namespace_service_accounts = ["external-dns:external-dns"]
    }
  }

}


output "dns_irsa_role_arn" {
  value = module.cert_manager_irsa_role.iam_role_arn
}
output "external_dns_irsa_role_arn" {
  value = module.external_dns_irsa_role.iam_role_arn
}

output "secret_irsa_role_arn" {
  value = module.iam_secret-role.iam_role_arn
}
