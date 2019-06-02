resource "aws_db_instance" "papi-db" {
  allocated_storage       = 20
  storage_type            = "gp2"
  backup_retention_period = 1
  engine                  = "postgres"
  engine_version          = "11.2"
  instance_class          = "db.t2.micro"
  name                    = "${var.db-name}"
  port                    = "${var.db-port}"
  username                = "${var.db-user}"
  password                = "${var.db-pass}"
  publicly_accessible     = false
  skip_final_snapshot     = true
  vpc_security_group_ids  = ["${aws_security_group.db_allow_srv.id}"]
}
