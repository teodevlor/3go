pipeline {
    agent none

    environment {
        IMAGE_NAME = "octotechvn/gogogo-api"
        STACK_NAME = "gogogo"
        DOCKER_CREDENTIALS = "credentials-docker-hub-octotechvn"
    }

    stages {
        stage('DEV Pipeline') {
            when { branch 'dev' }

            agent { label 'node-server-vanchuyenxanh' }

            stages {

                stage('Build & Push DEV Image') {
                    steps {
                        script {

                            def commit = sh(
                                script: "git rev-parse --short HEAD",
                                returnStdout: true
                            ).trim()

                            def buildNumber = env.BUILD_NUMBER
                            env.TAG = "dev-${buildNumber}-${commit}"

                            echo "Generated DEV TAG: ${env.TAG}"

                            withCredentials([usernamePassword(
                                credentialsId: env.DOCKER_CREDENTIALS,
                                usernameVariable: 'DOCKER_USER',
                                passwordVariable: 'DOCKER_PASS'
                            )]) {

                                sh """
                                    echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin

                                    docker build \
                                        -f docker/Dockerfile \
                                        -t ${IMAGE_NAME}:${TAG} .

                                    docker tag ${IMAGE_NAME}:${TAG} ${IMAGE_NAME}:dev-latest

                                    docker push ${IMAGE_NAME}:${TAG}
                                    docker push ${IMAGE_NAME}:dev-latest

                                    docker logout
                                """
                            }

                            currentBuild.description = "DEV Image: ${TAG}"
                        }
                    }
                }
            }
        }

        stage('PROD Pipeline') {
            when { branch 'master' }
            agent { label 'node-manager-docker-swam' }
            stages {
                stage('Build & Push PROD Image') {
                    steps {
                        script {

                            def commit = sh(
                                script: "git rev-parse --short HEAD",
                                returnStdout: true
                            ).trim()

                            def buildNumber = env.BUILD_NUMBER
                            env.TAG = "prod-${buildNumber}-${commit}"

                            echo "Generated PROD TAG: ${env.TAG}"

                            withCredentials([usernamePassword(
                                credentialsId: env.DOCKER_CREDENTIALS,
                                usernameVariable: 'DOCKER_USER',
                                passwordVariable: 'DOCKER_PASS'
                            )]) {

                                sh """
                                    echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin

                                    docker build \
                                        -f docker/Dockerfile \
                                        -t ${IMAGE_NAME}:${TAG} .

                                    docker tag ${IMAGE_NAME}:${TAG} ${IMAGE_NAME}:prod-latest

                                    docker push ${IMAGE_NAME}:${TAG}
                                    docker push ${IMAGE_NAME}:prod-latest

                                    docker logout
                                """
                            }

                            currentBuild.description = "PROD Image: ${TAG}"
                        }
                    }
                }

               stage('Run migrations and seeds') {
                    steps {
                        withCredentials([string(
                            credentialsId: 'prod-postgres-dsn',
                            variable: 'POSTGRES_DSN'
                        )]) {

                            sh """
                                set -ex

                                echo "=============================="
                                echo "DEBUG DATABASE CONNECTION"
                                echo "POSTGRES_DSN=\$POSTGRES_DSN"
                                echo "=============================="

                                echo "Running migrations..."

                                docker run --rm \
                                --network gogogo_gogogo_net \
                                -e POSTGRES_DSN="\$POSTGRES_DSN" \
                                ${IMAGE_NAME}:${TAG} \
                                ./migrate

                                echo "Running seeds..."

                                docker run --rm \
                                --network gogogo_gogogo_net \
                                -e POSTGRES_DSN="\$POSTGRES_DSN" \
                                ${IMAGE_NAME}:${TAG} \
                                ./migrate seed
                            """
                        }
                    }
                }

                stage('Deploy to Docker Swarm') {
                    steps {
                        sh """
                            export VERSION=${TAG}
                            docker stack deploy \
                                -c docker/docker-stack.yml \
                                ${STACK_NAME}
                        """
                    }
                }
            }
        }
    }

    post {
        success {
            echo "Pipeline completed successfully!"
        }
        failure {
            echo "Pipeline failed!"
        }
    }
}