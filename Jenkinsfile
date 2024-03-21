pipeline {

    agent any

    tools {
        go 'go1.14'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
    }

    stages {
        stage("Unit Test") {
            steps {
                echo 'UNIT TEST EXECUTION STARTED'
                sh 'make unit-tests'
            }
        }
        stage("Functional Test") {
            steps {
                echo 'FUNCTIONAL TEST EXECUTION STARTED'
                sh 'make functional-tests'
            }
        }
        stage("Build Image") {
            steps {
                echo 'BUILD EXECUTION STARTED'
                sh 'go version'
                sh 'go get ./...'
                sh 'docker build . -t satriagitu/product-go-micro'
            }
        }
        stage('Docker Push') {
            agent any
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', passwordVariable: 'dockerhubPassword', usernameVariable: 'dockerhubUser')]) {
                sh "docker login -u ${env.dockerhubUser} -p ${env.dockerhubPassword}"
                sh 'docker push satriagitu/product-go-micro'
                }
            }
        }
        stage('Deploy App on K8S with Helm') {
            steps {
                withCredentials([
                    string(credentialsId: 'my_kubernetes', variable: 'api_token')
                ]){
                    sh 'helm upgrade --install product-go-micro helm-product-go-micro -f helm-product-go-micro/values.yaml -n default --namespace=default  --kube-context=minikube'   
                }
            }
        }
    }
}
