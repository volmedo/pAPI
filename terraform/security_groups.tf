resource "aws_security_group" "srv_allow_ssh" {
  name        = "srv_allow_ssh"
  description = "Allow SSH inbound traffic to application server"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "srv_allow_http" {
  name        = "srv_allow_http"
  description = "Allow HTTP inbound traffic to application server"

  ingress {
    from_port   = "${var.srv-port}"
    to_port     = "${var.srv-port}"
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "srv_allow_outbound" {
  name        = "srv_allow_outbound"
  description = "Allow all outbound traffic from application server"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "db_allow_srv" {
  name        = "db_allow_srv"
  description = "Allow connections to the DB from the application server"

  ingress {
    from_port       = "${var.db-port}"
    to_port         = "${var.db-port}"
    protocol        = "tcp"
    security_groups = ["${aws_security_group.srv_allow_http.id}"]
  }
}
