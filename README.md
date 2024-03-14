# Golang Server Project Best Practices

### PS. Each Directory is a microservice and has its own Readme and each directory uses some new tech. So, do check it out.

## Dependency Injection :-
In simple words, we want our functions and packages to receive the objects they depend on ie. We dont want to declare new instances inside packages to have control over them. For Eg :- Using Structs to declare the private methods and save logger as a variables. The methods can access the value of logger using `g.logger` in their domain.

## Handle Timeouts :-
This is to prevent DOS attacks. Dont make requests make infinitely if your server crashes.

## Graceful Shutdown :-
Wait until current requests are handled and then shutdown the server. We use signal interrupts for this with channels

## Using JSON Encoder :-
Sometimes it is better to use Encoder than json.Marshal as we dont have to use additional buffer. This will matter when there are a lot of concurrent go routines being processed. It is also a bit faster.

## Validation of Data is very Important :-
Middleware ensures Connectivity between 2 or more types applications or components. You can write your validation code in middleware to make sure data is validated and then only goes to your handlers

## Running using Docker Image
`docker build . -t satriagitu/product-go-micro`

`docker run --network host -d <image-name/id>`

## Using CORS
Cross Origin Resource Sharing. Good security measure to protect the websites from malicious calls. It defines origins allowed to talk to the API. If source is not allowed, we reject the request.If we need to pass authentication headers like cookies the Origin Source should NOT be *. Like in this case where authentication using JWT is happening.

## File Handling
To Handle files usign a Golang server, we should not store them on our disk but on some cloud storage facility like S3. We can use our code to retrieve and send them. Golang's HTTP Fileserver helps to deal with sending the files from server. It provides a Handler for this.

To deal with very large files to make the data transfer efficient, we use GZipping. Time to compress and decompress is less than time taken for the unzipped file to be transferred.

## Note Deploy to Minikube (WSL)
sudo usermod -aG docker jenkins
sudo mkdir -p /var/lib/jenkins/.kube
sudo cp /path/to/your/kubeconfig /var/lib/jenkins/.kube/config 
sudo chown jenkins:jenkins /var/lib/jenkins/.kube/config
sudo systemctl restart jenkins
sudo chmod -R +r /home/satriagitu/.minikube/

sudo cp /home/satriagitu/.minikube/profiles/minikube/client.key /var/lib/jenkins/.minikube/profiles/minikube/client.key
sudo cp /home/satriagitu/.minikube/profiles/minikube/client.crt /var/lib/jenkins/.minikube/profiles/minikube/client.crt

update file config (sudo nano /var/lib/jenkins/.kube/config)
update location file client-certificate & client-key:
    client-certificate: /var/lib/jenkins/.minikube/profiles/minikube/client.crt
    client-key: /var/lib/jenkins/.minikube/profiles/minikube/client.key

## **(/path/to/your/kubeconfig = /home/satriagitu/.kube)

## Add file to github repo
    Add Jenkinsfile
    Add product-go-micro.yaml

## Configure Credential Kubernetes in Jenkins
    Configure kubernetes with secret token
    Configure cloud kubernetes with base64 minikube ca.crt

## Reference
    https://medium.com/@devayanthakur/minikube-configure-jenkins-kubernetes-plugin-25eb804d0dec
    https://medium.com/@vpkhandare2000/setting-up-jenkins-for-dockerized-minikube-deployment-c28d76e779b2