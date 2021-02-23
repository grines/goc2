resource "aws_instance" "goc2" {
  tags = {
    Name = "goc2-server",
    role = "goc2-server"
  }

  ami                         = "ami-011e27968f706be25"
  instance_type               = "t2.micro"
  iam_instance_profile        = "AmazonSSMRoleForInstancesQuickSetup"
  security_groups             = [aws_security_group.goc2_servers.id]
  subnet_id                   = aws_subnet.private.id
  associate_public_ip_address = false

  user_data = <<EOS
#!/bin/bash
cd /home/ubuntu
chmod +x goc2
./goc2 --web
EOS

  ebs_block_device {
    device_name = "/dev/sda1"
    volume_type = "standard"
    volume_size = 30
  }

  lifecycle {
    ignore_changes = [
      tags, security_groups, ebs_block_device
    ]
  }
}

resource "aws_instance" "redirector" {
  tags = {
    Name = "redirector",
    role = "redirector"
  }

  ami                         = "ami-04264ced02c5fea4d"
  instance_type               = "t2.micro"
  iam_instance_profile        = "AmazonSSMRoleForInstancesQuickSetup"
  security_groups             = [aws_security_group.redirectors.id]
  subnet_id                   = aws_subnet.public.id
  associate_public_ip_address = true

  user_data = <<EOS
#!/bin/bash
mkdir /var/log/redirector
chown -R ubuntu /var/log/redirector
EOS

  lifecycle {
    ignore_changes = [
      tags, security_groups
    ]
  }
}