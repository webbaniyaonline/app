pipeline {

```
agent any

environment {
    APP_NAME = "golang-app"
    CONTAINER_NAME = "golang-app"
}

stages {

    stage('Clone') {
        steps {
            git branch: 'main',
            url: 'https://github.com/webbaniyaonline/app.git'
        }
    }

    stage('Build Docker Image') {
        steps {
            sh 'docker build -t $APP_NAME .'
        }
    }

    stage('Stop Existing Container') {
        steps {
            sh '''
            docker stop $CONTAINER_NAME || true
            docker rm $CONTAINER_NAME || true
            '''
        }
    }

    stage('Run Container') {
        steps {
            sh '''
            docker run -d \
            --name $CONTAINER_NAME \
            -p 3000:3000 \
            --restart unless-stopped \
            $APP_NAME
            '''
        }
    }

}

post {
    success {
        echo 'Deployment Successful'
    }

    failure {
        echo 'Deployment Failed'
    }
}
```

}
