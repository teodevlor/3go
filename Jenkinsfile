pipeline {
    agent { label "${params.server}" }

    parameters {
        choice(
            name: 'server',
            choices: [
                'node-manager-docker-swam',
                'node-server-octotech-internal',
                'node-server-vanchuyenxanh'
            ]
        )

        choice(
            name: 'action',
            choices: ['deploy', 'rollback']
        )
    }

    stages {
        stage('Check VPS Info') {
            steps {
                sh 'whoami'
            }
        }
    }
}
