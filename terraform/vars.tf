variable "aws_region" {
  description = "The AWS region to deploy the resources"
  default     = "us-east-1"
}

variable "projectName" {
  description = "The name of the project"
  default     = "tc4-lanchonete-customer"
}

variable "mongo_user" {
  description = "MongoDB root username"
  type        = string
  sensitive   = true
}

variable "mongo_password" {
  description = "MongoDB root password"
  type        = string
  sensitive   = true
}

variable "mongo_db_name" {
  description = "MongoDB database name"
  type        = string
  default     = "tc4customer"
}
