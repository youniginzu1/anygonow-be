NAMESPACE=handyman
project=apiservice
kubectl create ns ${NAMESPACE} 
helm template deploy/production/${project} --values deploy/production/${project}/values.yaml
kubectl -n ${NAMESPACE} create configmap env --from-env-file deploy/dev/.env
kubectl -n ${NAMESPACE} get configmap
kubectl -n ${NAMESPACE} describe configmap env
helm install -n ${NAMESPACE} ${project} deploy/production/${project}
helm upgrade -n ${NAMESPACE} ${project} deploy/production/${project}
kubectl -n ${NAMESPACE} get pods
kubectl -n ${NAMESPACE} get svc
kubectl -n ${NAMESPACE} get deployment
kubectl -n ${NAMESPACE} describe deployment ${project}
helm uninstall -n ${NAMESPACE} ${project}
kubectl -n ${NAMESPACE} port-forward apiservice-9586d6785-9f78g 50051:50051