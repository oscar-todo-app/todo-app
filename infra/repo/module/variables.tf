variable "name" {
  description = "Name of the ECR Repository- should match the Github repo name."
  type        = string
  default     = "todo-app"
}

variable "oidc_arn" {
  description = "The OpenID Connect provider ARN"
  type        = string
}

variable "organization" {
  default = "oscar-todo-app"
}
