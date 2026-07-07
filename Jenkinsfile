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
                      -t sre-watchdog:${BUILD_NUMBER} \
                      .
                '''
            }
        }
    }

    post {
        success {
            echo 'SRE Watchdog pipeline completed successfully'
        }

        failure {
            echo 'SRE Watchdog pipeline failed'
        }

        always {
            echo "Build number: ${BUILD_NUMBER}"
        }
    }
}