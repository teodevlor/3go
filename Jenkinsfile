pipeline {
    agent none
    environment {
        IMAGE_NAME = "octotechvn/gogogo-api"
        STACK_NAME = "gogogo"
        DOCKER_CREDENTIALS = "credentials-docker-hub-octotechvn"
    }

    stages {

        // =========================
        // DEV PIPELINE
        // =========================
        stage('DEV Pipeline') {
            when {
                branch 'dev'
            }
            agent {
                label 'node-server-vanchuyenxanh'
            }
            stages {

                stage('Debug DEV') {
                    steps {
                        echo "Running on DEV node..."
                        echo "Branch: ${env.BRANCH_NAME}"
                        sh "whoami"
                        sh "hostname"
                    }
                }

                stage('Deploy DEV (temporary)') {
                    steps {
                        echo "Hiện tại DEV chưa build production image."
                    }
                }
            }
        }

        // =========================
        // PRODUCTION PIPELINE
        // =========================
        stage('PROD Pipeline') {
            when {
                branch 'master'
            }
            agent {
                label 'node-manager-docker-swam'
            }
            stages {

                stage('Debug PROD') {
                    steps {
                        echo "Running on PROD Swarm Manager..."
                        sh "whoami"
                        sh "hostname"
                        sh "git rev-parse --short HEAD"
                    }
                }

                stage('Build & Push Docker Image') {
                    steps {
                        script {

                            def commit = sh(
                                script: "git rev-parse --short HEAD",
                                returnStdout: true
                            ).trim()

                            def now = new Date().format(
                                "yyyyMMdd-HHmmss",
                                TimeZone.getTimeZone('Asia/Ho_Chi_Minh')
                            )

                            env.TAG = "${now}-${commit}"

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

                                    docker push ${IMAGE_NAME}:${TAG}

                                    docker logout
                                """
                            }

                            currentBuild.description = "Prod Image: ${TAG}"
                        }
                    }
                }

                stage('Run migrations and seeds') {
                    steps {
                        withCredentials([string(credentialsId: 'prod-postgres-dsn', variable: 'POSTGRES_DSN')]) {
                            sh """
                                set -e

                                echo "Running migrations..."

                                docker run --rm \
                                --network gogogo_gogogo_net \
                                -e POSTGRES_DSN="$POSTGRES_DSN" \
                                ${IMAGE_NAME}:${TAG} \
                                ./migrate

                                echo "Running seeds..."

                                docker run --rm \
                                --network gogogo_gogogo_net \
                                -e POSTGRES_DSN="$POSTGRES_DSN" \
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