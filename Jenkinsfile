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
            sh 'make report'
            cobertura autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'coverage.xml', conditionalCoverageTargets: '50, 0, 0', enableNewApi: true, failNoReports: false, failUnhealthy: false, failUnstable: false, lineCoverageTargets: '50, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '50, 0, 0', onlyStable: false
        }

    }
}
