node {
    // git credentialsId: 'private-key', url: 'git@github.com:by-sabbir/django-docker-jenkins.git'
    checkout scm
    try {

        stage ('Initializing Docker Repo') {
            withEnv(["DOCKER_USER=${DOCKER_USER}",
                     "DOCKER_PASSWORD=${DOCKER_PASSWORD}"]) {
                sh "make login" // sh "curl -X POST -H 'Content-Type: application/json' -d '{\"chat_id\": \"-1001556850823\", \"text\": \"Started ${JOB_BASE_NAME}_${JOB_NAME}- ${BUILD_URL}\", \"disable_notification\": false}' https://api.telegram.org/bot1750146504:AAE5lT-GQNVtEF48xQwH3IvecZa8WrytYY8/sendMessage"
            }
        }
        stage ('Unit Test') {
            sh "make test"
        }
        stage ('Build Image and Publish'){
            sh "make publish"
        }

        stage ("Deploying") {
            ansiblePlaybook colorized: true, credentialsId: 'private-docker-creds-id', inventory: 'ansible/hosts', playbook: 'ansible/playbook/rollout.yaml'
        }
        
    }
    finally {
        stage ("Cleaning Up..."){
            sh 'make cleanup'
            sh 'make logout'
        }

        stage ("report") {
            sh 'make reporthtml'
            publishHTML (target : [allowMissing: true,
                alwaysLinkToLastBuild: true,
                keepAll: true,
                reportDir: 'reports',
                reportFiles: 'index.html',
                reportName: 'UnitTest Report',
                reportTitles: 'Unit Tests'])
        }

    }
}
