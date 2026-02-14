.PHONY: help deploy-dev deploy-prod destroy-dev destroy-prod status-dev status-prod

NAMESPACE_DEV = efact-dev
NAMESPACE_PROD = efact

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-25s %s\n", $$1, $$2}'

deploy-dev: ## Deploy to development
	kubectl apply -k overlays/dev/
	@$(MAKE) status-dev

deploy-prod: ## Deploy to production
	kubectl apply -k overlays/prod/
	@$(MAKE) status-prod

destroy-dev: ## Destroy development deployment
	kubectl delete -k overlays/dev/ --ignore-not-found=true

destroy-prod: ## Destroy production deployment
	kubectl delete -k overlays/prod/ --ignore-not-found=true

status-dev: ## Show development status
	@kubectl get pods,svc,pvc -n $(NAMESPACE_DEV)

status-prod: ## Show production status
	@kubectl get pods,svc,pvc -n $(NAMESPACE_PROD)

logs-ms1-dev: ## Logs MS1 dev
	kubectl logs -n $(NAMESPACE_DEV) -l app=ms1-documents -f --tail=100

logs-ms1-prod: ## Logs MS1 prod
	kubectl logs -n $(NAMESPACE_PROD) -l app=ms1-documents -f --tail=100

logs-ms2-dev: ## Logs MS2 dev
	kubectl logs -n $(NAMESPACE_DEV) -l app=ms2-validator -f --tail=100

logs-ms2-prod: ## Logs MS2 prod
	kubectl logs -n $(NAMESPACE_PROD) -l app=ms2-validator -f --tail=100

port-forward-dev: ## Port-forward dev services
	@kubectl port-forward -n $(NAMESPACE_DEV) svc/ms1-documents-service 5000:5000 & \
	kubectl port-forward -n $(NAMESPACE_DEV) svc/ms2-validator-service 8080:8080 & \
	kubectl port-forward -n $(NAMESPACE_DEV) svc/rabbitmq-service 15672:15672 & \
	kubectl port-forward -n $(NAMESPACE_DEV) svc/kibana-service 5601:5601 & \
	wait

port-forward-prod: ## Port-forward prod services
	@kubectl port-forward -n $(NAMESPACE_PROD) svc/ms1-documents-service 5000:5000 & \
	kubectl port-forward -n $(NAMESPACE_PROD) svc/ms2-validator-service 8080:8080 & \
	kubectl port-forward -n $(NAMESPACE_PROD) svc/rabbitmq-service 15672:15672 & \
	kubectl port-forward -n $(NAMESPACE_PROD) svc/kibana-service 5601:5601 & \
	wait

scale-ms1-dev: ## Scale MS1 dev (REPLICAS=N)
	kubectl scale deployment ms1-documents -n $(NAMESPACE_DEV) --replicas=$(REPLICAS)

scale-ms2-dev: ## Scale MS2 dev (REPLICAS=N)
	kubectl scale deployment ms2-validator -n $(NAMESPACE_DEV) --replicas=$(REPLICAS)

scale-ms1-prod: ## Scale MS1 prod (REPLICAS=N)
	kubectl scale deployment ms1-documents -n $(NAMESPACE_PROD) --replicas=$(REPLICAS)

scale-ms2-prod: ## Scale MS2 prod (REPLICAS=N)
	kubectl scale deployment ms2-validator -n $(NAMESPACE_PROD) --replicas=$(REPLICAS)

restart-ms1-dev: ## Restart MS1 dev
	kubectl rollout restart deployment ms1-documents -n $(NAMESPACE_DEV)

restart-ms2-dev: ## Restart MS2 dev
	kubectl rollout restart deployment ms2-validator -n $(NAMESPACE_DEV)

restart-ms1-prod: ## Restart MS1 prod
	kubectl rollout restart deployment ms1-documents -n $(NAMESPACE_PROD)

restart-ms2-prod: ## Restart MS2 prod
	kubectl rollout restart deployment ms2-validator -n $(NAMESPACE_PROD)

hpa-dev: ## Create HPA dev
	kubectl autoscale deployment ms1-documents -n $(NAMESPACE_DEV) --cpu-percent=70 --min=1 --max=5
	kubectl autoscale deployment ms2-validator -n $(NAMESPACE_DEV) --cpu-percent=70 --min=1 --max=5

hpa-prod: ## Create HPA prod
	kubectl autoscale deployment ms1-documents -n $(NAMESPACE_PROD) --cpu-percent=70 --min=2 --max=10
	kubectl autoscale deployment ms2-validator -n $(NAMESPACE_PROD) --cpu-percent=70 --min=2 --max=10

clean-dev: ## Delete dev namespace
	kubectl delete namespace $(NAMESPACE_DEV) --ignore-not-found=true

clean-prod: ## Delete prod namespace
	kubectl delete namespace $(NAMESPACE_PROD) --ignore-not-found=true
