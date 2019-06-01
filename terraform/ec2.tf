resource "aws_instance" "papi-server" {
  ami           = "${data.aws_ami.amazon-linux-2.id}"
  instance_type = "t2.micro"
  key_name      = "${aws_key_pair.ssh-key.key_name}"

  vpc_security_group_ids = [
    "${aws_security_group.srv_allow_ssh.id}",
    "${aws_security_group.srv_allow_http.id}",
    "${aws_security_group.srv_allow_outbound.id}",
  ]

  connection {
    type        = "ssh"
    user        = "ec2-user"
    private_key = "${file("${var.ssh-key-path}")}"
  }

  provisioner "file" {
    source      = "${var.srv-bin-path}"
    destination = "/home/ec2-user/${basename(var.srv-bin-path)}"
  }

  provisioner "file" {
    source      = "${var.db-migrations-path}"
    destination = "/home/ec2-user/"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /home/ec2-user/${basename(var.srv-bin-path)}",
      <<EOF
      nohup /home/ec2-user/${basename(var.srv-bin-path)} \
        -port=${var.srv-port} \
        -dbhost=${aws_db_instance.papi-db.address} \
        -dbport=${var.db-port} \
        -dbuser=${var.db-user} \
        -dbpass=${var.db-pass} \
        -dbname=${var.db-name} \
        -migrations="/home/ec2-user/${basename(var.db-migrations-path)}" &
      EOF
      ,
      "sleep 1",
    ] # this feels hacky :S
  }
}

data "aws_ami" "amazon-linux-2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

resource "aws_key_pair" "ssh-key" {
  key_name   = "ssh-key"
  public_key = "${file("${var.ssh-key-path}.pub")}"
}

output "srv-ip" {
  value = "${aws_instance.papi-server.public_ip}"
}

output "srv-port" {
  value = "${var.srv-port}"
}
