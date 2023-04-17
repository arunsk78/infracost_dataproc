terraform {
  backend "gcs" {
    bucket  = "my-foo-bucket-tfstate"
    prefix  = "dataproc"
  }

  required_version = ">= 0.12"
}
provider "google-beta" {
  project = "my-foo-project"
  region  = "europe-west3"
  zone = "europe-west3-c"
}

module "my_foo_cluster" {
  source = "./modules/terraform-google-dataproc"
  cluster_name = "my-cool-cluster"
  cluster_version = "1.4"
  region = "europe-west3"
  master_ha = false
  zone = "europe-west3-c"
  master_instance_type = "n1-standard-4"
  service_account = "my-cool-account@my-cool-project.iam.gserviceaccount.com"
  network = "my-cool-network"
  worker_instance_type = "n1-standard-4"
  conda_packages = "pandas=0.23.4 scikit-learn=0.20.0 pytest=3.8.0 pyyaml=3.13"
  pip_packages = "gensim==3.7.1 logdecorator==2.1"
  staging_bucket = "my-cool-bucket"
}