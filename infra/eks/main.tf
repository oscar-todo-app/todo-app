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


data "aws_vpc" "default" {
cidr_block = "10.0.0.0/16"
tags = {
name = "todo-vpc"
}
}
resource "random_string" "todo-db-password" {
  length  = 32
  upper   = true
  numeric = true
  special = false
}
resource "aws_security_group" "todo-db" {
  vpc_id      = "${data.aws_vpc.default.id}"
  name        = "todo-db"
  description = "Allow all inbound for Postgres"
  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
resource "aws_db_instance" "todo-db" {
  identifier             = "todo-db"
  instance_class         = "db.t4g.micro"
  allocated_storage      = 5
  engine                 = "postgres"
  engine_version         = "16.1"
  skip_final_snapshot    = true
  publicly_accessible    = true
  vpc_security_group_ids = [aws_security_group.todo-db.id]
  username               = "todo"
  password               = "random_string.todo-db-password.result}"}
