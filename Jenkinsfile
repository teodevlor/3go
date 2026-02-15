pipeline {
    agent { label 'gogogo-api' }

    environment {
        APP_NAME = "gogogo-demo"
    }

    stages {

        stage('Checkout Source') {
            steps {
                echo "🔄 Checking out source..."
                checkout scm
            }
        }

        stage('Verify VPS') {
            steps {
                echo "🖥 Verifying deploy server..."
                sh '''
                    echo "Current user:"
                    whoami
                    echo "Hostname:"
                    hostname
                    echo "Working directory:"
                    pwd
                    echo "IP address:"
                    hostname -I
                '''
            }
        }

        stage('Git Info') {
            steps {
                echo "📦 Git information"
                sh '''
                    echo "Branch: ${BRANCH_NAME}"
                    echo "Commit SHA:"
                    git rev-parse --short HEAD
                    echo "Last commit message:"
                    git log -1 --pretty=%B
                '''
            }
        }

        stage('Build (Demo)') {
            steps {
                echo "🏗 Building project..."
                sh '''
                    echo "Simulating build..."
                    sleep 2
                    echo "Build completed"
                '''
            }
        }

        stage('Deploy (Demo)') {
            steps {
                echo "🚀 Deploying..."
                sh '''
                    echo "Deploying ${APP_NAME} on $(hostname)"
                    echo "Deploy success ✅"
                '''
            }
        }
    }

    post {
        success {
            echo "🎉 Pipeline completed successfully!"
        }
        failure {
            echo "❌ Pipeline failed!"
        }
    }
}
