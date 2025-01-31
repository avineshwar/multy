resource "aws_s3_object" "file1_public_aws" {
  provider     = "aws.eu-west-1"
  bucket       = aws_s3_bucket.obj_storage_aws.id
  key          = "index.html"
  acl          = "public-read"
  content      = "${templatefile("index.html", { "aws_url" = "${aws_api_gateway_deployment.lambda_aws.invoke_url}", "azure_url" = "https://${azurerm_function_app.lambda_azure.default_hostname}/trigger" })}"
  content_type = "text/html"
}
data "archive_file" "lambda_aws" {
  type        = "zip"
  source_dir  = "source_dir/aws"
  output_path = ".multy/tmp/multy_function_aws.zip"
}
resource "aws_lambda_function" "lambda_aws" {
  provider         = "aws.eu-west-1"
  tags             = { "Name" = "multy_function" }
  function_name    = "multy_function"
  role             = aws_iam_role.iam_for_lambda_lambda_aws.arn
  filename         = ".multy/tmp/multy_function_aws.zip"
  source_code_hash = data.archive_file.lambda_aws.output_base64sha256
  runtime          = "python3.9"
  handler          = "lambda_function.lambda_handler"
}
resource "aws_iam_role" "iam_for_lambda_lambda_aws" {
  provider           = "aws.eu-west-1"
  tags               = { "Name" = "iam_for_lambda_lambda_aws" }
  name               = "iam_for_lambda_lambda_aws"
  assume_role_policy = "{\"Statement\":[{\"Action\":[\"sts:AssumeRole\"],\"Effect\":\"Allow\",\"Principal\":{\"Service\":\"lambda.amazonaws.com\"}}],\"Version\":\"2012-10-17\"}"
}
resource "aws_iam_role_policy_attachment" "lambda_aws" {
  provider   = "aws.eu-west-1"
  role       = aws_iam_role.iam_for_lambda_lambda_aws.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
resource "aws_api_gateway_rest_api" "lambda_aws" {
  provider    = "aws.eu-west-1"
  tags        = { "Name" = "multy_function" }
  name        = "multy_function"
  description = ""
}
resource "aws_api_gateway_resource" "lambda_proxy" {
  provider    = "aws.eu-west-1"
  rest_api_id = aws_api_gateway_rest_api.lambda_aws.id
  parent_id   = aws_api_gateway_rest_api.lambda_aws.root_resource_id
  path_part   = "{proxy+}"
}
resource "aws_api_gateway_method" "lambda_proxy" {
  provider      = "aws.eu-west-1"
  rest_api_id   = aws_api_gateway_rest_api.lambda_aws.id
  resource_id   = aws_api_gateway_resource.lambda_proxy.id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "lambda_proxy" {
  provider                = "aws.eu-west-1"
  rest_api_id             = aws_api_gateway_rest_api.lambda_aws.id
  resource_id             = aws_api_gateway_method.lambda_proxy.resource_id
  http_method             = aws_api_gateway_method.lambda_proxy.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_aws.invoke_arn
}
resource "aws_api_gateway_method" "lambda_proxy_root" {
  provider      = "aws.eu-west-1"
  rest_api_id   = aws_api_gateway_rest_api.lambda_aws.id
  resource_id   = aws_api_gateway_rest_api.lambda_aws.root_resource_id
  http_method   = "ANY"
  authorization = "NONE"
}
resource "aws_api_gateway_integration" "lambda_proxy_root" {
  provider                = "aws.eu-west-1"
  rest_api_id             = aws_api_gateway_rest_api.lambda_aws.id
  resource_id             = aws_api_gateway_method.lambda_proxy_root.resource_id
  http_method             = aws_api_gateway_method.lambda_proxy_root.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambda_aws.invoke_arn
}
resource "aws_api_gateway_deployment" "lambda_aws" {
  provider    = "aws.eu-west-1"
  rest_api_id = aws_api_gateway_rest_api.lambda_aws.id
  stage_name  = "api"
  depends_on  = [aws_api_gateway_integration.lambda_proxy, aws_api_gateway_integration.lambda_proxy_root]
}
resource "aws_lambda_permission" "lambda_aws" {
  provider      = "aws.eu-west-1"
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "multy_function"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_rest_api.lambda_aws.execution_arn}/*/*"
}
resource "aws_s3_bucket" "obj_storage_aws" {
  provider = "aws.eu-west-1"
  bucket   = "test-storage-9999919"
}
resource "azurerm_resource_group" "fun-rg" {
  name     = "fun-rg"
  location = "northeurope"
}
data "archive_file" "lambda_azure" {
  type        = "zip"
  source_dir  = "source_dir/azure"
  output_path = ".multy/tmp/multy_function_azure.zip"
}
resource "azurerm_storage_account" "lambda_azure" {
  resource_group_name             = azurerm_resource_group.fun-rg.name
  name                            = "multyfunctionstacwzob"
  location                        = "northeurope"
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  allow_nested_items_to_be_public = false
  blob_properties {
    versioning_enabled = false
  }
}
resource "azurerm_function_app" "lambda_azure" {
  resource_group_name        = azurerm_resource_group.fun-rg.name
  name                       = "multyfunction"
  location                   = "northeurope"
  storage_account_name       = azurerm_storage_account.lambda_azure.name
  storage_account_access_key = azurerm_storage_account.lambda_azure.primary_access_key
  app_service_plan_id        = azurerm_app_service_plan.lambda_azure.id
  os_type                    = "linux"
  provisioner "local-exec" {
    command = "az functionapp deployment source config-zip -g ${self.resource_group_name} -n ${self.name} --src ${data.archive_file.lambda_azure.output_path}"
  }
}
resource "azurerm_app_service_plan" "lambda_azure" {
  resource_group_name = azurerm_resource_group.fun-rg.name
  name                = "multyfunctionsvplj8si"
  location            = "northeurope"
  kind                = "Linux"
  reserved            = true
  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
provider "aws" {
  region = "eu-west-1"
  alias  = "eu-west-1"
}


provider "azurerm" {
  features {
  }
}
