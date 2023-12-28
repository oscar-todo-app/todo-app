terraform {

  required_providers {
    helm = {
      source = "hashicorp/helm"
    }
  }
}
provider "helm" {
  kubernetes {
    host                   = var.clusterHost
    cluster_ca_certificate = base64decode(var.clusterToken)
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      args        = ["eks", "get-token", "--cluster-name", var.clusterName]
      command     = "aws"
    }
  }
}
resource "helm_release" "argocd-install" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  create_namespace = true
  timeout          = 600
}
