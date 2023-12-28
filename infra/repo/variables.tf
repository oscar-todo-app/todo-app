variable "name" {
  description = "Name of the ECR Repository- should match the Github repo name."
  type        = string
}

variable "organization" {
  description = "Name of the Github Organization."
  type        = string
  default     = "multi-py"
}
