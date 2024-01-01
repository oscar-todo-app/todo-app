terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.31.0"
    }
  }
}
provider "aws" {
  region = "eu-west-2"
}

locals {
  cluster_name = "todo-cluster"
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.4.0"

  name                 = "todo-vpc"
  cidr                 = "172.16.0.0/16"
  azs                  = ["eu-west-2a", "eu-west-2b", "eu-west-2c"]
  private_subnets      = ["172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"]
  public_subnets       = ["172.16.4.0/24", "172.16.5.0/24", "172.16.6.0/24"]
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
  public_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                      = "1"
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${local.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"             = "1"
  }
}
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name    = local.cluster_name
  cluster_version = "1.28"

  enable_irsa                    = true
  cluster_endpoint_public_access = true

  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
  }
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets
  eks_managed_node_groups = {
    general = {
      desired_size = 1
      min_size     = 1
      max_size     = 1

      labels = {
        role = "general"
      }

      instance_types = ["t3a.xlarge"]
      capacity_type  = "ON_DEMAND"
    }
  }
}

module "db" {
  source     = "./db/"
  vpcID      = module.vpc.vpc_id
  secGroupID = module.eks.cluster_security_group_id
  subnets    = module.vpc.private_subnets
}


module "argo" {
  source       = "./argo/"
  clusterHost  = module.eks.cluster_endpoint
  clusterName  = module.eks.cluster_name
  clusterToken = module.eks.cluster_certificate_authority_data
}

module "iam" {
  source       = "./iam/"
  secret_arn   = module.db.secret_arn
  provider_url = module.eks.cluster_oidc_issuer_url
  provider_arn = module.eks.oidc_provider_arn
}


output "cert-manager" {
  value = module.iam.dns_irsa_role_arn
}

output "external-dns" {
  value = module.iam.external_dns_irsa_role_arn
}

output "secrets-manager" {
  value = module.iam.secret_irsa_role_arn
}
output "db_endpoint" {
  value = module.db.instance_endpoint

}
