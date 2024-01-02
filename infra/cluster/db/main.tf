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


resource "random_password" "todo-pass" {
  length           = 16
  special          = true
  override_special = "_!%@^"

}

resource "aws_db_subnet_group" "todo-db-subnet-group" {
  name       = "todo-db-subnet-group"
  subnet_ids = var.subnets
}

resource "aws_secretsmanager_secret" "db-pass" {
  recovery_window_in_days = 0
  name                    = "db-pass-new"
}

resource "aws_secretsmanager_secret_version" "db-pass-v" {
  secret_id     = aws_secretsmanager_secret.db-pass.id
  secret_string = random_password.todo-pass.result
}


resource "aws_security_group" "todo-db-group" {
  name        = "todo-db-group"
  description = "Allow postgress"
  vpc_id      = var.vpcID
}

resource "aws_security_group_rule" "todo-db" {
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 5432
  to_port                  = 5432
  description              = "Allow all inbound for Postgres"
  security_group_id        = aws_security_group.todo-db-group.id
  source_security_group_id = var.secGroupID
}
resource "aws_db_instance" "todo-db" {
  identifier             = "todo-db"
  instance_class         = "db.t4g.micro"
  allocated_storage      = 5
  engine                 = "postgres"
  engine_version         = "16.1"
  skip_final_snapshot    = true
  publicly_accessible    = false
  multi_az               = false
  vpc_security_group_ids = [aws_security_group.todo-db-group.id]
  username               = "todo"
  password               = random_password.todo-pass.result
  db_name                = "todo"
  db_subnet_group_name   = aws_db_subnet_group.todo-db-subnet-group.name
}


output "secret_arn" {
  value = aws_secretsmanager_secret.db-pass.arn
}

output "instance_endpoint" {
  value = aws_db_instance.todo-db.endpoint
}
