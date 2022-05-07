resource "azurerm_resource_group" "rg1" {
  name     = "rg1"
  location = "northeurope"
}
provider "azurerm" {
  features {}
}
provider "aws" {
  region = "eu-west-1"
  alias  = "eu-west-1"
}
