pipeline {
    agent any

    environment {
        // Сhange the required parameters
        DOCKER_HUB_USER = 'antifootbolist'
        REPO_URL='https://github.com/antifootbolist/store-api.git'
        GHP_URL='https://github.com/antifootbolist/antifootbolist.github.io.git'
        GIT_AUTHOR_NAME = "Alexey Borodulin"
        GIT_AUTHOR_EMAIL = "antifootbolist@gmail.com"
        
        // Create variable in Jenkins
        // PROD_IP - IP address of PROD server where we deploy 'main' branch
        // STAGE_IP - IP address of STAGE server where we deploy 'develop' branch
        
        // Configure on Jenkins and change
        GH_TOKEN_ID='antifootbolist-github-access-token'
        DOCKER_HUB_LOGIN='docker_hub_login'

        // Optional to change
        GO_APP_PORT = '8080'
        GO_APP_NAME = 'go-app'
        NGINX_PORT = '80'
        NGINX_NAME = 'nginx'
        PG_PORT = '5432'
        PG_NAME = 'postgresql'
        APIDOC_NAME = 'apidoc'
        APP_NET = 'app-net'
    }

    stages {
        stage('Build') {
            steps {
                script {
                    // Define variables depending on branch name
                    serverIp = env.BRANCH_NAME == 'main' ? env.PROD_IP : env.STAGE_IP
                    testData = env.BRANCH_NAME == 'main' ? 'False' : 'True'

                    // Save the variables as environment variables
                    env.SERVER_IP = serverIp
                    env.TEST_DATA = testData
                    
                    def app_names = [env.PG_NAME, env.GO_APP_NAME, env.NGINX_NAME]
                    for (app_name in app_names) {
                        app = docker.build("${DOCKER_HUB_USER}/${app_name}", "-f ${app_name}/Dockerfile .")
                        docker.withRegistry('https://registry.hub.docker.com', env.DOCKER_HUB_LOGIN) {
                            app.push("${env.BUILD_NUMBER}")
                            app.push("latest")
                        }
                    }
                }
            }
        }
        stage ('Deploy DB') {
            steps {
                script {
                    sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker pull ${DOCKER_HUB_USER}/${env.PG_NAME}:${env.BUILD_NUMBER}\""
                    try {
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker stop ${env.PG_NAME}\""
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker rm ${env.PG_NAME}\""
                    } catch (err) {
                        echo: 'caught error: $err'
                    }
                    sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker run -d --restart always --name ${env.PG_NAME} --network ${APP_NET} -p ${env.PG_PORT}:${env.PG_PORT} ${DOCKER_HUB_USER}/${env.PG_NAME}:${env.BUILD_NUMBER}\""
                }
            }
        }
        stage ('Migrate DB schema') {
            steps {
                script {
                    sh "cat flyway/flyway.conf" 
                    sh "sed \"s/localhost/${SERVER_IP}/\" flyway/flyway.conf > flyway/flyway.conf"
                    sh "cat flyway/flyway.conf"
                    sh "docker run --rm -v flyway/sql:/flyway/sql -v flyway/conf:/flyway/conf flyway/flyway:9.8.1 migrate"
                }
            }
        }
        stage ('Deploy App') {
            steps {
                script {
                    sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker pull ${DOCKER_HUB_USER}/${env.GO_APP_NAME}:${env.BUILD_NUMBER}\""
                    try {
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker stop ${env.GO_APP_NAME}\""
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker rm ${env.GO_APP_NAME}\""
                    } catch (err) {
                        echo: 'caught error: $err'
                    }
                    sh "scp -o StrictHostKeyChecking=no ${env.GO_APP_NAME}/env.list ${SERVER_IP}:./env.list"
                    sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker run -d --restart always --name ${env.GO_APP_NAME} --network ${APP_NET} -p ${env.GO_APP_PORT}:${env.GO_APP_PORT} --env-file ./env.list -e TEST_DATA=${TEST_DATA} ${DOCKER_HUB_USER}/${env.GO_APP_NAME}:${env.BUILD_NUMBER}\""
                }
            }
        }
        stage ('Deploy Nginx') {
            steps {
                script {
                    sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker pull ${DOCKER_HUB_USER}/${env.NGINX_NAME}:${env.BUILD_NUMBER}\""
                    try {
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker stop ${env.NGINX_NAME}\""
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker rm ${env.NGINX_NAME}\""
                    } catch (err) {
                        echo: 'caught error: $err'
                    }
                    sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker run -d --restart always --name ${env.NGINX_NAME} --network ${APP_NET} -p ${env.NGINX_PORT}:${env.NGINX_PORT} ${DOCKER_HUB_USER}/${env.NGINX_NAME}:${env.BUILD_NUMBER}\""
                }
            }
        }
        /*
        stage ('Deploy') {
            steps {
                script {
                    // Setting variables depending on the branch
                    if (env.BRANCH_NAME == 'main') {
                        SERVER_IP = env.PROD_IP
                        TEST_DATA = 'False'
                    } else {
                        SERVER_IP = env.STAGE_IP
                        TEST_DATA = 'True'
                    }
                    
                    def app_names = [env.PG_NAME, env.FLYWAY_NAME, env.GO_APP_NAME, env.NGINX_NAME]
                    for (app_name in app_names) {
                        if (app_name == env.PG_NAME) {
                            sh 'echo "Deploying PostgreSQL server"'
                            app_port = env.PG_PORT
                        }
                        if (app_name == env.FLYWAY_NAME) {
                            sh 'echo "Running scheme migration"'
                        }
                        if (app_name == env.GO_APP_NAME) {
                            sh 'echo "Deploying Go application"'
                            app_port = env.GO_APP_PORT
                            sh "scp -o StrictHostKeyChecking=no ${app_name}/env.list ${SERVER_IP}:./env.list"
                        }
                        if (app_name == env.NGINX_NAME) {
                            sh 'echo "Deploying Nginx Web server"'
                            app_port = env.NGINX_PORT
                        }
                        sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker pull ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                        try {
                            sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker stop ${app_name}\""
                            sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker rm ${app_name}\""
                        } catch (err) {
                            echo: 'caught error: $err'
                        }
                        if (app_name == env.GO_APP_NAME) {
                            sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker run -d --restart always --name ${app_name} --network ${APP_NET} -p ${app_port}:${app_port} --env-file ./env.list -e TEST_DATA=${TEST_DATA} ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                        } else {
                            if (app_name == env.FLYWAY_NAME) {
                                sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker run --rm --name ${app_name} --network ${APP_NET} ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER} migrate\"" 
                            } else {
                                sh "ssh -o StrictHostKeyChecking=no ${SERVER_IP} \"docker run -d --restart always --name ${app_name} --network ${APP_NET} -p ${app_port}:${app_port} ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                            }
                        }
                    }
                }
            }
        }
        */
        stage('apiDoc') {
            when {
                branch 'main'
            }
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: env.GH_TOKEN_ID, usernameVariable: 'GH_NAME', passwordVariable: 'GH_TOKEN')]) {
                        def url = sh(script: "echo ${env.GHP_URL} | sed 's#https://##'", returnStdout: true).trim()
                        
                        def app = docker.build("${env.APIDOC_NAME}", "-f ${env.GO_APP_NAME}/Dockerfile.apidoc .")
                        
                        app.inside("--workdir=/app -u root") {
                            sh "git clone https://${GH_NAME}:${GH_TOKEN}@${url} /app/ghp_repo"
                            sh 'cp -R /app/apidoc/* /app/ghp_repo/'
                            try {
                                sh "cd /app/ghp_repo && \
                                git config user.name ${GIT_AUTHOR_NAME} && \
                                git config user.email ${GIT_AUTHOR_EMAIL} && \
                                git add . && \
                                git commit -m \"Update apiDoc for build #${env.BUILD_NUMBER}\" && \
                                git push -f https://${GH_NAME}:${GH_TOKEN}@${url} main"
                            }
                            catch (err) {
                                echo: 'Skipping apiDoc update. Reason: $err'
                            } 
                        }
                    }
                }
            }
        }
    }
    post {
        always {
            script {
                cleanWs()
            }
        }
    }
}
