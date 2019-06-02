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

variable "db-name" {
  description = "Name of the DB to create in the DB server instance"
}

variable "db-port" {
  description = "Port where the DB server will accept connections"
  default     = 5432
}

variable "db-user" {
  description = "Username to use when accessing the DB"
}

variable "db-pass" {
  description = "Password to use when accessing the DB"
}

variable "db-migrations-path" {
  description = "Path to the directory that contains the migration files to be deployed on provisioning"
}
