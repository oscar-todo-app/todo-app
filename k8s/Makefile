setup-cluser:
	aws eks --region eu-west-2 update-kubeconfig --name todo-cluster
get-pass:
	kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
apply-sets:
	kubectl apply -f ./applicatonset.yaml
	kubectl apply -f ./todoset.yaml.yaml

apply-ingress:
	kubectl apply -f ./ingress.yaml
all: setup-cluser get-pass apply-sets apply-ingress
