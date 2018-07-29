pipeline {

  agent {
    label {
      label "master"
    }
  }
  options {
    // gitLabConnection('gitlab')
    buildDiscarder(logRotator(numToKeepStr: '5'))
    disableConcurrentBuilds()
    timestamps()
  }
  triggers {
    // gitlab(triggerOnPush: true, triggerOnMergeRequest: true, branchFilterType: 'All')
    pollSCM '@hourly'
  }
  stages {
    stage('Build API') {
      steps {
        sh 'docker image build --no-cache \
            --tag bcgdv/api:latest --build-arg GIT_COMMIT=$(git log -1 --format=%H) .'
      }
    }
    
    stage('Deploy API') {
      steps {
      sh 'API_HOST_PORT=5000 docker stack deploy --compose-file docker-compose.yml airQuality'
      }
    }
  }
}
