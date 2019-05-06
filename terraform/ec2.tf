resource "aws_instance" "papi-server" {
#  ami           = "ami-09ead922c1dad67e4"
  ami = "${data.aws_ami.amazon-linux-2.id}"
  instance_type = "t2.micro"
}

data "aws_ami" "amazon-linux-2" {
  most_recent = true
  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-ebs"]
  }
}

output "papi-server-ip" {
  value = "${aws_instance.papi-server.public_ip}"
}
