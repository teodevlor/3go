pipeline {
    agent any

    stages {

        stage('Debug') {
            steps {
                echo "Branch: ${env.BRANCH_NAME}"
                sh "whoami"
                sh "hostname"
            }
        }

        stage('Deploy Dev') {
            when {
                branch 'dev'
            }
            steps {
                echo "Deploying DEV..."
                echo "Hẹ hẹ..."
                echo "Hẹ hẹ hẹ gẹ gẹ..."
            }
        }

        stage('Deploy Prod') {
            when {
                branch 'master'
            }
            steps {
                echo "Deploying PRODUCTION..."
            }
        }
    }
}