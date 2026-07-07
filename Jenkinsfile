pipeline {
    agent any

    options {
        timestamps()
        skipDefaultCheckout(true)
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Set Image Tag') {
            steps {
                script {
                    env.IMAGE_TAG = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                }

                echo "Image tag: ${env.IMAGE_TAG}"
            }
        }

        stage('Go Test') {
            steps {
                sh '''
                    docker run --rm \
                      -v "$WORKSPACE":/workspace \
                      -w /workspace \
                      golang:alpine \
                      go test ./...
                '''
            }
        }

        stage('Go Vet') {
            steps {
                sh '''
                    docker run --rm \
                      -v "$WORKSPACE":/workspace \
                      -w /workspace \
                      golang:alpine \
                      go vet ./...
                '''
            }
        }

        stage('Docker Build') {
            steps {
                sh '''
                    docker build \
                      -t sre-watchdog:${IMAGE_TAG} \
                      .
                '''
            }
        }
    }

    post {
        success {
            echo "Pipeline completed successfully"
            echo "Built image: sre-watchdog:${IMAGE_TAG}"
        }

        failure {
            echo 'Pipeline failed'
        }
    }
}