pipeline {
    agent any

    environment {
        // Need to change
        PROD_IP = '130.193.36.79'
        GO_APP_PORT = '8080'
        GO_APP_NAME = 'go-app'
        NGINX_PORT = '80'
        NGINX_NAME = 'nginx'
        PG_PORT = '5432'
        PG_NAME = 'postgresql'
        APIDOC_NAME = 'apidoc'
        APP_NET = 'app-net'
        DOCKER_HUB_USER = 'antifootbolist'
        REPO_URL='https://github.com/antifootbolist/store-api.git'
        GHP_URL='https://github.com/antifootbolist/antifootbolist.github.io.git'
        GH_TOKEN_ID='antifootbolist-github-access-token'
        GIT_AUTHOR_NAME = "Alexey Borodulin"
        GIT_AUTHOR_EMAIL = "antifootbolist@gmail.com"
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
                        // Don't forget to create docker_hub_login credential to autorize on Docker Hub
                        docker.withRegistry('https://registry.hub.docker.com', 'docker_hub_login') {
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
                    def app = docker.build("${env.APIDOC_NAME}", "-f ${env.GO_APP_NAME}/Dockerfile.apidoc .")

                    app.inside {
                        sh 'cp -R /app/apidoc ./apidoc'

                        withCredentials([usernamePassword(credentialsId: env.GH_TOKEN_ID, usernameVariable: 'GH_NAME', passwordVariable: 'GH_TOKEN')]) {
                            def url = sh(script: "echo ${env.GHP_URL} | sed 's#https://##'", returnStdout: true).trim()

                            sh 'pwd'
                            sh "git clone https://${GH_NAME}:${GH_TOKEN}@${url} ghp_repo"
                            sh 'cp -R ./apidoc/* ghp_repo/'
                            sh "git config user.name ${GIT_AUTHOR_NAME}"
                            sh "git config user.email ${GIT_AUTHOR_EMAIL}"
                            sh 'pwd'
                            sh 'cd ghp_repo && git add .'
                            sh 'pwd'
                            sh 'cd ghp_repo && git commit -m "Update apiDoc documentation for Store API"'
                            sh 'pwd'
                            sh "cd ghp_repo && git push -f https://${GH_NAME}:${GH_TOKEN}@${url} main"
                            sh 'pwd'
                        }
                    }
                }
            }
        }
        stage ('Deploy Apps to Prod') {
            steps {
                // Don't forget to create prod_login credential to autorize on Prod server
                withCredentials ([usernamePassword(credentialsId: 'prod_login', usernameVariable: 'USERNAME', passwordVariable: 'USERPASS')]) {
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
                            }
                            if (app_name == env.NGINX_NAME) {
                                sh 'echo "Deploying Nginx Web server"'
                                app_port = env.NGINX_PORT
                            }
                            sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$PROD_IP \"docker pull ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
                            try {
                                sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$PROD_IP \"docker stop ${app_name}\""
                                sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$PROD_IP \"docker rm ${app_name}\""
                            } catch (err) {
                                echo: 'caught error: $err'
                            }
                            sh "sshpass -p '$USERPASS' -v ssh -o StrictHostKeyChecking=no $USERNAME@$PROD_IP \"docker run -d --restart always --name ${app_name} --network ${APP_NET} -p ${app_port}:${app_port} ${DOCKER_HUB_USER}/${app_name}:${env.BUILD_NUMBER}\""
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
