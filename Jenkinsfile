pipeline {

  agent any
  tools{go 'Go'}
  
    environment {
        // Define Quay.io credentials for pushing images
        QUAY_IO_CREDENTIALS = 'quay-io-credentials'
        QUAY_IO_REGISTRY = 'https://quay.io'
        QUAY_IO_USERNAME = 'udaykirank'
        QUAY_IO_REPOSITORY_PREFIX = 'csye7125group3/controller' // Customize this as needed
        IMAGE_NAME = 'quay.io/csye7125group3/controller'
    }

  stages {

    stage('Checkout') {
      steps {
        git branch: 'main', 
            credentialsId: 'github-token-jenkins',
            url: 'https://github.com/cyse7125-fall2023-group03/kube-operator.git'
      }
    }

    // stage('Build Image') {
    //   steps {  
    //     script {
    //       def latestVersion = sh(returnStdout: true, script: "git describe --tags --abbrev=0").trim()
    //       docker.build("${env.IMAGE_NAME}:${latestVersion}", '.')
    //       docker.build("${env.IMAGE_NAME}:latest", '.')
    //     }
    //   }
    // }
    
    stage('Login to Quay.io') {
        steps {
            withCredentials([usernamePassword(credentialsId: QUAY_IO_CREDENTIALS, passwordVariable: 'QUAY_IO_PASSWORD', usernameVariable: 'QUAY_IO_USERNAME')]) {
                script {
                        sh "docker login -u $QUAY_IO_USERNAME -p $QUAY_IO_PASSWORD ${env.QUAY_IO_REGISTRY}"
                        sh "make docker-build docker-push IMG=quay.io/csye7125group3/controller:$latestVersion"
                        sh "make docker-build docker-push IMG=quay.io/csye7125group3/controller:latest"
                    }
                }
        }
    }
    
    // stage('Push Docker Image to Quay.io') {
    //     steps {
    //         script {
    //             def latestVersion = sh(returnStdout: true, script: "git describe --tags --abbrev=0").trim()
    //             docker.withRegistry("https://quay.io/", QUAY_IO_CREDENTIALS) {
    //                 docker.image("${env.IMAGE_NAME}:${latestVersion}").push()
    //                 docker.image("${env.IMAGE_NAME}:latest").push()
    //             }
    //         }
    //     }
    // }

  }

}
