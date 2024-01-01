resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["a031c46782e6e6c662c2c87c76da9aa62ccabd8e"]
}

module "repository" {
  source   = "./module/"
  name     = "todo-app"
  oidc_arn = aws_iam_openid_connect_provider.github.arn
}

output "role_arn" {
  value = module.repository.ecr_role
}
