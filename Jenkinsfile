pipeline {
  agent {
    node {
      label 'master'
    }

  }
  stages {
    stage('Build API') {
      steps {
        sh 'docker image build --no-cache --tag bcgdv/api:latest --build-arg GIT_COMMIT=$(git log -1 --format=%H) .'
      }
    }
    stage('Deploy API') {
      steps {
        sh 'API_HOST_PORT=5000 docker stack deploy --compose-file docker-compose.yml airQuality'
      }
    }
    stage('Test API') {
      steps {
        sh 'docker exec -it airQuality_api.1.$(docker service ps -f \'name=airQuality_api.1\' airQuality_api -q --no-trunc | head -n1) go test -v'
      }
    }
  }
  options {
    buildDiscarder(logRotator(numToKeepStr: '5'))
    disableConcurrentBuilds()
    timestamps()
  }
  triggers {
    pollSCM('@hourly')
  }
}
