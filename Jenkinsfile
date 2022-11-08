node {
    checkout scm
    try {

        stage ('Initializing Docker Repo') {
            withEnv(["DOCKER_USER=${DOCKER_USER}",
                     "DOCKER_PASSWORD=${DOCKER_PASSWORD}"]) {
                sh "make login"
            }
        }
        stage ('Unit Test') {
            sh "make test"
        }
        stage ('Build Image and Publish'){
            sh "make publish"
        }

        stage ("Deploying") {
            shortCommit = sh(returnStdout: true, script: "git log -n 1 --pretty=format:'%h'").trim()
            echo "GIT_COMMIT is ${shortCommit}"
            ansiblePlaybook extras: "-e commitId=${shortCommit}", inventory: 'ansible/hosts', playbook: 'ansible/playbook/rollout.yaml'
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
