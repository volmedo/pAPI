resource "aws_instance" "papi-server" {
  ami           = "${data.aws_ami.amazon-linux-2.id}"
  instance_type = "t2.micro"
  key_name      = "${aws_key_pair.ssh-key.key_name}"

  vpc_security_group_ids = [
    "${aws_security_group.allow_ssh.id}",
    "${aws_security_group.allow_http.id}",
    "${aws_security_group.allow_outbound.id}",
  ]

  connection {
    type        = "ssh"
    user        = "ec2-user"
    private_key = "${file("${var.ssh-key-path}")}"
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

output "papi-server-ip" {
  value = "${aws_instance.papi-server.public_ip}"
}
