pipeline {
agent any

```
environment {
    APP_NAME = "golang-app"
    CONTAINER_NAME = "golang-app"
    APP_PORT = "3000"
}

stages {

    stage('Checkout') {
        steps {
            echo 'Checking out source code...'
            checkout scm
        }
    }

    stage('Build Docker Image') {
        steps {
            echo 'Building Docker image...'
            sh '''
                docker build -t ${APP_NAME}:latest .
            '''
        }
    }

    stage('Stop Existing Container') {
        steps {
            echo 'Stopping old container if exists...'
            sh '''
                docker stop ${CONTAINER_NAME} || true
                docker rm ${CONTAINER_NAME} || true
            '''
        }
    }

    stage('Deploy Container') {
        steps {
            echo 'Deploying application...'
            sh '''
                docker run -d \
                --name ${CONTAINER_NAME} \
                --restart unless-stopped \
                -p ${APP_PORT}:${APP_PORT} \
                ${APP_NAME}:latest
            '''
        }
    }

    stage('Verify Deployment') {
        steps {
            echo 'Checking running containers...'
            sh '''
                docker ps
            '''
        }
    }
}

post {
    success {
        echo 'Application deployed successfully.'
    }

    failure {
        echo 'Deployment failed.'
    }

    always {
        echo 'Pipeline execution completed.'
    }
}
```

}
