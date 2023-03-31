base_image:=gcr.io/cloudrun/hello
load_envs:=. ./.env

all: init apply

.PHONY: init
init:
	${load_envs} && terraform init

.PHONY: plan
plan:
	${load_envs} && terraform plan

.PHONY: apply
apply:
	${load_envs} && terraform apply

.PHONY: destroy
destroy:
	${load_envs} && terraform destroy
