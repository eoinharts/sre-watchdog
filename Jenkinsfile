pipeline {
    agent any

    options {
        timestamps()
        skipDefaultCheckout(true)
    }

    parameters {
        booleanParam(
            name: 'PUBLISH_TO_ECR',
            defaultValue: false,
            description: 'Push the built image to AWS ECR'
        )
    }

    environment {
        AWS_REGION = 'eu-west-1'
        ECR_REPOSITORY = 'sre-watchdog'
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

        stage('Set ECR Metadata') {
            when {
                expression {
                    params.PUBLISH_TO_ECR
                }
            }

            steps {
                script {
                    env.AWS_ACCOUNT_ID = sh(
                        script: '''
                            aws sts get-caller-identity \
                              --query Account \
                              --output text
                        ''',
                        returnStdout: true
                    ).trim()

                    env.ECR_REGISTRY =
                        "${env.AWS_ACCOUNT_ID}.dkr.ecr.${env.AWS_REGION}.amazonaws.com"

                    env.ECR_IMAGE =
                        "${env.ECR_REGISTRY}/${env.ECR_REPOSITORY}:${env.IMAGE_TAG}"
                }

                echo "ECR image: ${env.ECR_IMAGE}"
            }
        }

        stage('ECR Login') {
            when {
                expression {
                    params.PUBLISH_TO_ECR
                }
            }

            steps {
                sh '''
                    aws ecr get-login-password \
                      --region "${AWS_REGION}" \
                    | docker login \
                      --username AWS \
                      --password-stdin "${ECR_REGISTRY}"
                '''
            }
        }

        stage('Tag Image for ECR') {
            when {
                expression {
                    params.PUBLISH_TO_ECR
                }
            }

            steps {
                sh '''
                    docker tag \
                      sre-watchdog:${IMAGE_TAG} \
                      "${ECR_IMAGE}"
                '''
            }
        }

        stage('Push Image to ECR') {
            when {
                expression {
                    params.PUBLISH_TO_ECR
                }
            }

            steps {
                sh '''
                    docker push "${ECR_IMAGE}"
                '''
            }
        }
    }

    post {
        success {
            echo 'SRE Watchdog pipeline completed successfully'

            script {
                if (params.PUBLISH_TO_ECR) {
                    echo "Published image: ${env.ECR_IMAGE}"
                } else {
                    echo "Built local image: sre-watchdog:${env.IMAGE_TAG}"
                }
            }
        }

        failure {
            echo 'SRE Watchdog pipeline failed'
        }

        always {
            echo "Build number: ${env.BUILD_NUMBER}"

            script {
                if (params.PUBLISH_TO_ECR && env.ECR_REGISTRY) {
                    sh '''
                        docker logout "${ECR_REGISTRY}" || true
                    '''
                }
            }
        }
    }
}