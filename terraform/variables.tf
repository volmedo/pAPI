variable "ssh-key-path" {
  description = "Path to the key to use to SSH into deployed instances for provisioning"
}

variable "srv-bin-path" {
  description = "Path to the server binary executable that will be deployed and run on provisioning"
}

variable "srv-port" {
  description = "Port where the server will accept requests"
  default     = 8080
}
