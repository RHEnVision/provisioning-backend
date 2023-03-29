##@ Dashboard

.PHONY: generate-dashboard
generate-dashboard: ## Generate dashboard
	cp ./deploy/grafana-dashboard-provisioning.tmpl.yaml ./deploy/dashboards/grafana-dashboard-provisioning.configmap.yaml
	json_pp < ./deploy/Provisioning-1674548289785.json | sed 's/^/    /' >> ./deploy/dashboards/grafana-dashboard-provisioning.configmap.yaml

.PHONY: validate-dashboard
validate-dashboard: generate-dashboard ## Compare dashboard configmaps with git
	git diff --exit-code deploy/dashboards/*.configmap.yaml
