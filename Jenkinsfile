pipeline {
    agent any

    environment {
        APP_NAME = "gogogo-demo"
    }

    stages {

        stage('Checkout') {
            steps {
                echo "🔄 Checking out source..."
                checkout scm
            }
        }

        stage('Info') {
            steps {
                echo "📦 Branch: ${env.BRANCH_NAME}"
                sh 'echo "Commit SHA: $(git rev-parse --short HEAD)"'
                sh 'echo "Last commit message:"'
                sh 'git log -1 --pretty=%B'
            }
        }

        stage('Build') {
            steps {
                echo "🏗 Building project..."
                sh 'echo "Simulating build..."'
                sh 'sleep 3'
            }
        }

        stage('Test') {
            steps {
                echo "🧪 Running tests..."
                sh 'echo "All tests passed ✅"'
            }
        }

        stage('Deploy (Demo)') {
            steps {
                echo "🚀 Deploying (fake deploy)..."
                sh 'echo "Deploy success"'
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
