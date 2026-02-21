// ======================
// GIT
// ======================
def cloneAndCheckout(String commit, String githubUrl) {

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
            git init
            git remote add origin ${githubUrl}

            git config credential.helper store
            echo "https://\$GIT_USER:\$GIT_PASS@github.com" > ~/.git-credentials

            git fetch origin
            git checkout ${commit}
        """
    }
}


// ======================
// DOCKER BUILD
// ======================
def buildAndPushImage(String imageName, String commit) {

    def now = new Date().format("yyyyMMdd-HHmmss", TimeZone.getTimeZone('Asia/Ho_Chi_Minh'))
    def tag = "${now}-${commit}"

    withCredentials([usernamePassword(
        credentialsId: 'credentials-docker-hub-octotechvn',
        usernameVariable: 'DOCKER_USER',
        passwordVariable: 'DOCKER_PASS'
    )]) {

        sh """
            echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin
            docker build -f docker/Dockerfile -t ${imageName}:${tag} .
            docker push ${imageName}:${tag}
            docker logout
        """
    }

    return tag
}


// ======================
// DEPLOY
// ======================
def deployStack(String tag) {

    if (!tag) {
        error("Docker tag is required!")
    }

    sh """
        export VERSION=${tag}
        docker stack deploy -c docker/docker-stack.yml gogogo
    """
}


// ======================
// DEPLOY FLOW
// ======================
def deployFlow(String commit, String githubUrl, String imageName) {

    stage('Clone & Checkout') {
        cloneAndCheckout(commit, githubUrl)
    }

    def tag

    stage('Build & Push Docker Image') {
        tag = buildAndPushImage(imageName, commit)
        currentBuild.description = "Image tag: ${tag}"
    }

    stage('Deploy Stack to Swarm') {
        deployStack(tag)
    }
}


// ======================
// ROLLBACK FLOW
// ======================
def rollbackFlow(String selectedImage) {

    if (!selectedImage) {
        error("Docker image tag is required for rollback!")
    }

    def tag = selectedImage.contains(":") ?
              selectedImage.split(":")[-1] :
              selectedImage

    stage('Rollback Stack') {
        deployStack(tag)
    }
}


// ======================
// MAIN PIPELINE
// ======================
node(params.server) {

    def IMAGE_NAME = "octotechvn/gogogo-api"
    def GITHUB_URL = "https://github.com/teodevlor/3go.git"

    def ACTION = params.action
    def COMMIT = params.GIT_COMMIT_ID?.trim()
    def DOCKER_IMAGE = params.DOCKER_IMAGE?.trim()

    currentBuild.displayName = "${ACTION}"

    try {

        if (ACTION == 'deploy') {

            if (!COMMIT) {
                error("Commit ID is required for deploy!")
            }

            deployFlow(COMMIT, GITHUB_URL, IMAGE_NAME)
        }

        if (ACTION == 'rollback') {
            rollbackFlow(DOCKER_IMAGE)
        }

        stage('Success') {
            echo "Pipeline completed successfully!"
        }

    } catch (err) {

        stage('Failure') {
            echo "Pipeline failed: ${err}"
        }

        throw err
    }
}
