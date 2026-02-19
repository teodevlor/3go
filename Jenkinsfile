def cloneAndCheckout(String commit) {

    if (!commit) {
        error("Commit ID is empty!")
    }

    deleteDir()

    withCredentials([usernamePassword(
        credentialsId: 'github-teodevlor',
        usernameVariable: 'GIT_USER',
        passwordVariable: 'GIT_PASS'
    )]) {

        sh """
            echo "📦 Cloning repository..."

            git init
            git remote add origin https://github.com/teodevlor/3go.git

            git config credential.helper store
            echo "https://\$GIT_USER:\$GIT_PASS@github.com" > ~/.git-credentials

            git fetch origin
            git checkout ${commit}

            echo "✅ Checkout completed"
        """
    }
}


def buildAndPushImage(String imageName, String commit) {

    withCredentials([usernamePassword(
        credentialsId: 'credentials-docker-hub-octotechvn',
        usernameVariable: 'DOCKER_USER',
        passwordVariable: 'DOCKER_PASS'
    )]) {

        sh """
            echo "🔐 Logging in to Docker Hub..."
            echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin

            echo "🏗 Building Docker image ${imageName}:${commit}..."
            docker build -f docker/Dockerfile -t ${imageName}:${commit} .

            echo "📤 Pushing image..."
            docker push ${imageName}:${commit}

            docker logout

            echo "✅ Image pushed successfully"
        """
    }
}


def deployStack(String commit) {

    sh """
        echo "🚀 Deploying Docker Stack..."

        export VERSION=${commit}

        docker stack deploy \
            -c docker-stack.yml \
            gogogo

        echo "📦 Current Stack Services:"
        docker stack services gogogo

        echo "📊 Stack Tasks:"
        docker stack ps gogogo

        echo "✅ Stack deployed successfully!"
    """
}


def rollbackStack(String commit) {

    sh """
        echo "↩ Rolling back stack to ${commit}..."

        export VERSION=${commit}

        docker stack deploy \
            -c docker-stack.yml \
            gogogo

        echo "✅ Rollback completed!"
    """
}


node(params.server) {

    def IMAGE_NAME = "octotechvn/gogogo-api"
    def COMMIT     = params.GIT_COMMIT_ID?.trim()
    def ACTION     = params.action

    stage('Build Info') {
        echo "====== BUILD INFO ======"
        echo "Server: ${params.server}"
        echo "Action: ${ACTION}"
        echo "Commit: ${COMMIT}"
        echo "========================"
    }

    try {

        if (ACTION == 'deploy') {

            stage('Clone & Checkout') {
                cloneAndCheckout(COMMIT)
            }

            stage('Build & Push Docker Image') {
                buildAndPushImage(IMAGE_NAME, COMMIT)
            }

            stage('Deploy Stack to Swarm') {
                deployStack(COMMIT)
            }
        }

        if (ACTION == 'rollback') {

            stage('Rollback Stack') {
                rollbackStack(COMMIT)
            }
        }

        stage('Success') {
            echo "🎉 Pipeline completed successfully!"
        }

    } catch (err) {

        stage('Failure') {
            echo "❌ Pipeline failed: ${err}"
        }

        throw err
    }
}
