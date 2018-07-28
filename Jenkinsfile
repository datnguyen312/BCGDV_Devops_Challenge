pipeline {
  agent any
  stages {
    stage('Build API') {
      steps {
        sh 'docker image build --no-cache \
            --tag bcgdv/api:latest --build-arg GIT_COMMIT=$(git log -1 --format=%H) .'
      }
    }
    stage('Deploy API') {
      steps {
        echo 'Deploying API'
      }
    }
  }
}
