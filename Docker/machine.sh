docker-machine rm default
docker-machine create --driver virtualbox \
  --virtualbox-hostonly-cidr "25.0.1.100/24" \
  --virtualbox-disk-size "5120" \
  --virtualbox-cpu-count "2" \
  --virtualbox-memory "4096" \
  --engine-insecure-registry 10.47.7.214 \
  default

docker-machine ssh default "sudo mkdir -p /etc/docker/certs.d/"
docker-machine ssh default "sudo cp /Users/amanpreet.singh/.minikube/files/etc/ssl/certs/* /etc/docker/certs.d/"

docker-machine ssh default "sudo mkdir -p  /var/lib/boot2docker/certs/"
docker-machine ssh default "sudo cp /Users/amanpreet.singh/.minikube/files/etc/ssl/certs/* /var/lib/boot2docker/certs/"

docker-machine restart