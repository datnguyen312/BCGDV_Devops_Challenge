# BCGDV Infrastructure Engineer take home test

## Air Quality API

You are on a venture team that is building a business in the environmental sector. As part of this effort, some engineers have built a small HTTP API in Golang that allows consumers to discover the air quality in their local area. The API receives latitude and longitude coordinates via a HTTP GET request, and subsequently retrieves air quality data for those coordinates by making a further HTTP request to the extneral [Air Visual API](https://airvisual.com/api) . Once the Air Visual response is received, our API returns a JSON response containing the nearest City to the requested coordindates along with the latest air quality index in that location.

- Air quality URL: `/air-quality?lat=50&lon=1`

To aid with analytics and future evolution of the API, each incoming query is stored in a PostreSQL database, and a second external HTTP GET endpoint exposes all the recently queried cities.

- Queried cities URL: `/queried-cities`

As the venture continue and the business evolves, the API's responsibilities will inevitably change. It is also unknown what the future levels of traffic for the API could be, as this is dependent on the success of the venture.

# The Challenge - Part 1

The team have requested your assistance and expertise to help deploy the Air Quality API and associated database to Amazon Web Services.

- Automate deployment of the Air Quality API to AWS
- You may use any tools or services you choose
- Make any code changes you feel are necessary to facilitate the deployment

Challenge tips:
- Think about how your solution will support future evolution and changes to the API
- We value simple solutions over anything complex. Don't overengineer
- This doesn't need to be production-ready code, but be prepared to explain what further changes you'd make if you had more time

AWS Access:
- The [AWS Free Tier](https://aws.amazon.com/free/) should provide you with enough free usage to complete this task. However, if you would like to use any services that fall outside the free tier, then we'll happily reimburse you for this usage (providing it's fair and you haven't just been mining Bitcoin)

# The Challenge - Part 2

The team have also requested your advice for their ongoing development and ownership of the Air Quality API.

Put together a short presentation that details:

- A proposed Continuous Delivery pipeline for the Air Quality API
- Changes that should be made to the API to facilitate its operation in production

Challenge tips:
- Consider security, resiliency, monitoring etc


## Initial Setup
- Install Docker locally from [https://docs.docker.com/install/](https://docs.docker.com/install/)
- Run `start.sh` to start the API
- Visit [http://localhost:5000](http://localhost:5000) in a browser to confirm the API is running
- Run `test.sh` to execute tests

# Changes made _(to the original code repository; .zip)_ by Shahzad Chaudhry
- Added a vagrantfile that creates two Ubuntu VMs which form a docker swarm mode cluster _(1x master and 1x worker)_ for testing API locally
  - Ubuntu 16.04.5 LTS (Xenial Xerus)
  - Docker version 18.06.0-ce _(as of 28/07/2018)_
- Updated .gitignore file to exclude Vagrant logs
- Moved Dockerfile outside of api folder
- Removed .dockerignore file as it is not needed
- Changed orchestration technology from being "docker compose" to "docker swarm mode"
- Added a Jenkins compose file. Jenkins will act as a build & deployment automation engine
- Made substential changes to docker-compose.yml
  - Changed compose file format from 3 to 3.6
  - Removed "links" as it is deprecated in the latest compose file format
  - No longer publishing DB port numbers as api service connects internally to db over a software defined network
  - Created a network over which all services talk to each other internally
  - Added docker secret management. And so, password from the compose file is removed
  - Added a service called "adminer". This is a web based full-featured database management tool written in PHP
- Updated start.sh which:
  - builds the api docker image and tags it as bcgdv/api:latest
  - Runs both api and db as two docker swarm mode services in a single stack
  - Shows an example how multiple stacks can be started. This means running multiple api versions at the same time
- Added "docker-compose.portainer.yml" that can be used to start [Portainer](https://portainer.io/) which is an open-source lightweight management UI which allows you to easily manage your Docker hosts or Swarm clusters.

## Part 1 answer
This section addresses how to deploy code changes in this repository to a runtime environment in AWS. This all needs to be automated in CI / CD fashion; push changes to this repo -> a Jenkins job is triggered automatically to build code -> resultant docker image is then deployed in in a runtime in AWS
- Firstly, an AWS infra will need to be setup. See ref no. 2 blow for a link on how to setup a docker swarm cluster using "Docker Community Edition (CE) for AWS". Here are some of the benefits for such a cluster:
  - Native to Docker
  - Skip the boilerplate and maintenance work
  - Minimal, Docker-focused base
  - Self-cleaning and self-healing
  - Logging native to the platforms
- Secondly, A Jenkins instance will need to be installed. This tool will help with build and deployment automation via declarative syntext; pipeline-as-a-code. See start.sh file for instructions on how to start jenkins
- Thirdly, create a webhook in github repo that on a push event will trigger a build in Jenkins. Webhook URL format = "YourJenkinsURL/github-webhook/""

# References
1. See https://github.com/shazChaudhry/vagrant repo for instructions on how to create Ubuntu VMs using Vagrant and how to start Portainer for local testing
1. Follow the instructions at https://docs.docker.com/docker-for-aws/ to setup a docker swarm cluster in AWS; 1x master and 1x work
