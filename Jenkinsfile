pipeline {
    agent any

    environment {
        // Сhange the required parameters     
        PROD_IP = '130.193.36.79'
        DOCKER_HUB_USER = 'antifootbolist'
        REPO_URL='https://github.com/antifootbolist/store-api.git'
        GHP_URL='https://github.com/antifootbolist/antifootbolist.github.io.git'
        GIT_AUTHOR_NAME = "Alexey Borodulin"
        GIT_AUTHOR_EMAIL = "antifootbolist@gmail.com"
        
        // Configure on Jenkins and change
        GH_TOKEN_ID='antifootbolist-github-access-token'
        DOCKER_HUB_LOGIN='docker_hub_login'
        PROD_LOGIN='prod_login'


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
        stage('Checkout') {
            steps {
                git url: env.REPO_URL,
                    branch: 'main'
            }
        }
        stage('Build and Push Docker Images') {
            steps {
                script {
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
        stage('Generate and Publish apiDoc') {
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
        stage ('Deploy API to Prod') {
            steps {
                withCredentials ([usernamePassword(credentialsId: env.PROD_LOGIN, usernameVariable: 'USERNAME', passwordVariable: 'USERPASS')]) {
                    script {
                        def app_names = [env.PG_NAME, env.GO_APP_NAME, env.NGINX_NAME]
                        for (app_name in app_names) {
                            if (app_name == env.PG_NAME) {
                                sh 'echo "Deploying PostgreSQL server"'
                                app_port = env.PG_PORT
                            }
                            if (app_name == env.GO_APP_NAME) {
                                sh 'echo "Deploying Go application"'
                                app_port = env.GO_APP_PORT
                                sh "sshpass -p '${USERPASS}' scp -o StrictHostKeyChecking=no ${app_name}/env.list ${USERNAME}@${PROD_IP}:/home/${USERNAME}/env.list"
                            }
                            if (app_name == env.NGINX_NAME) {
                                sh 'echo "Deploying Nginx Web server"'
                                app_port = env.NGINX_PORT
                            }
                            sh "sshpass -p '${USERPASS}' -v ssh -o StrictHostKeyChecking=no ${USERNAME}@${PROD_IP} \"docker pull ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                            try {
                                sh "sshpass -p '${USERPASS}' -v ssh -o StrictHostKeyChecking=no ${USERNAME}@${PROD_IP} \"docker stop ${app_name}\""
                                sh "sshpass -p '${USERPASS}' -v ssh -o StrictHostKeyChecking=no ${USERNAME}@${PROD_IP} \"docker rm ${app_name}\""
                            } catch (err) {
                                echo: 'caught error: $err'
                            }
                            if (app_name == env.GO_APP_NAME) {
                                sh "sshpass -p '${USERPASS}' -v ssh -o StrictHostKeyChecking=no ${USERNAME}@${PROD_IP} \"docker run -d --restart always --name ${app_name} --network ${APP_NET} -p ${app_port}:${app_port} --env-file /home/${USERNAME}/env.list ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                            } else {
                                sh "sshpass -p '${USERPASS}' -v ssh -o StrictHostKeyChecking=no ${USERNAME}@${PROD_IP} \"docker run -d --restart always --name ${app_name} --network ${APP_NET} -p ${app_port}:${app_port} ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                            }

                            // sh "sshpass -p '${USERPASS}' -v ssh -o StrictHostKeyChecking=no ${USERNAME}@${PROD_IP} \"docker run -d --restart always --name ${app_name} --network ${APP_NET} -p ${app_port}:${app_port} ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
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
