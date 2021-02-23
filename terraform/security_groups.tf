// GO C2 Servers

resource "aws_security_group" "goc2_servers" {
  name        = "goc2-servers"
  description = "Go C2 servers"

  vpc_id = aws_vpc.vpc.id
}

resource "aws_security_group_rule" "goc2_server_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.goc2_servers.id
}

resource "aws_security_group_rule" "goc2_server_ssh" {
  type                     = "ingress"
  from_port                = 22
  to_port                  = 22
  protocol                 = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.goc2_servers.id
}

resource "aws_security_group_rule" "goc2_servers_admin" {
  type                     = "ingress"
  from_port                = 8005
  to_port                  = 8005
  protocol                 = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.goc2_servers.id
}

// Redirectors

resource "aws_security_group" "redirectors" {
  name        = "redirectors"
  description = "Redirector"

  vpc_id = aws_vpc.vpc.id
}

resource "aws_security_group_rule" "redirector_egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.redirectors.id
}

resource "aws_security_group_rule" "redirectors_http" {
  type        = "ingress"
  from_port   = 80
  to_port     = 80
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.redirectors.id
}

resource "aws_security_group_rule" "redirectors_https" {
  type        = "ingress"
  from_port   = 443
  to_port     = 443
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.redirectors.id
}
