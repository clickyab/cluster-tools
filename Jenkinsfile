node {
    stage('Prepare') {
        checkout scm
    }
    stage('Build') {
        checkout scm
        sh "./bin/ci-test.sh all"
    }
    stage('Lint') {
        checkout scm
        sh "./bin/ci-test.sh lint"
    }
    stage('Test') {
        checkout scm
        sh "./bin/ci-test.sh test"
    }
}
