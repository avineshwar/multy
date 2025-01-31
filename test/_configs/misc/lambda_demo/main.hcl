multy "lambda" "lambda" {
  function_name   = "multy_function"
  runtime         = "python3.9"
  source_code_dir = cloud_specific_value({ aws : "source_dir/aws", azure : "source_dir/azure" })
}

multy "object_storage" "obj_storage" {
  clouds = ["aws"]
  name   = "test-storage-9999919"
}
multy "object_storage_object" "file1_public" {
  clouds         = ["aws"]
  name           = "index.html"
  content        = templatefile("index.html", { aws_url = aws.lambda.url, azure_url = "${azure.lambda.url}/trigger" })
  object_storage = obj_storage
  content_type   = "text/html"
  acl            = "public_read"
}

config {
  location = var.location
  clouds   = ["aws", "azure"]
}
variable "location" {
  type    = string
  default = "EU_WEST_1"
}