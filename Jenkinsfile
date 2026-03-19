pipeline {
    agent none

    environment {
        // 공통 환경 변수 정의
        IMAGE_REPO = 'repo.iris.tools/template'
        IMAGE_NAME = 'test'
        // Jenkins 가 실행되는 노드의 docker 그룹 ID
        DOCKER_GROUP_ID = '998'
    }

    stages {
        // 빌드 및 테스트 스테이지 (Go 이미지 사용)
        stage('CI Pipeline') {
            agent {
                docker {
                    image 'golang:1.26.1-alpine3.23'
                    // DooD(Docker Out of Docker) 환경으로 Jenkins 의 workspace 내
                    // Go 모듈 캐시 디렉토리를 Host 에 마운트하여 재사용한다.
                    // docker socket 사용 시 permission 문제 해결을 위해 group-add 옵션을 사용하여 docker 그룹에 추가한다.
                    args '''
                        -v /var/run/docker.sock:/var/run/docker.sock
                        -v /DATA/jenkins-2025/.go/pkg/mod:/root/go/pkg/mod
                        -v /DATA/jenkins-2025/.go/build-cache:/root/.cache/go-build
                        --group-add ${DOCKER_GROUP_ID}
                    '''
                    reuseNode true
                }
            }
            stages {
                stage('Checkout') {
                    steps {
                        echo '소스 코드 체크아웃 중...'
                        checkout scm
                    }
                }

                stage('Setup') {
                    steps {
                        echo '빌드 도구 설치 중...'
                        sh 'apk add --no-cache make build-base git'
                        sh 'mkdir -p build'
                    }
                }

                stage('Build') {
                    steps {
                        echo '프로젝트 빌드 중...'
                        sh 'make build'
                    }
                }

                stage('Lint') {
                    steps {
                        echo 'golangci-lint 실행 중...'
                        sh 'make lint'
                    }
                }

                stage('Test') {
                    steps {
                        echo '테스트 실행 중...'
                        sh 'go test ./... -v -covermode=count -coverprofile=build/coverage.out 2>&1 | tee build/test-output.txt'
                        sh 'go tool cover -html=build/coverage.out -o build/cov-out.html'
                    }
                    post {
                        always {
                            publishHTML([
                                reportDir: 'build',
                                reportFiles: 'cov-out.html',
                                reportName: 'Go Coverage Report',
                                keepAll: true,
                                alwaysLinkToLastBuild: true,
                                allowMissing: true
                            ])
                        }
                    }
                }

            }
        }

        // SonarQube 분석 스테이지 (sonar-scanner-cli 이미지 사용)
        stage('SonarQube Analysis') {
            agent {
                docker {
                    image 'sonarsource/sonar-scanner-cli:latest'
                    reuseNode true
                }
            }
            steps {
                echo 'SonarQube 분석 및 리포트 생성 중...'
                withSonarQubeEnv('sonar-iris-tools') {
                    // 프로젝트 설정은 sonar-project.properties 에서 자동으로 읽음
                    sh 'sonar-scanner'
                }
            }
        }

        // 도커 빌드 및 푸시 스테이지 (Jenkins 노드에서 직접 실행)
        // Jenkins 컨테이너에 docker CLI가 설치되어 있고 /var/run/docker.sock에 접근 가능해야 함
        stage('Docker Build & Push') {
            steps {
                script {
                    def fullImageName = "${env.IMAGE_REPO}/${env.IMAGE_NAME}"

                    // Jenkins 환경변수에서 Git Hash 가져오기
                    // checkout scm 단계에서 GIT_COMMIT 환경변수가 생성됨
                    def gitHash = env.GIT_COMMIT ? env.GIT_COMMIT.take(7) : 'dev'

                    echo "빌드할 이미지: ${fullImageName}"
                    echo "Git Hash: ${gitHash}"

                    // Docker 인증 및 빌드/푸시
                    // mobigen-harbor 는 인증 정보로 Jenkins 의 Credentials 에 등록되어 있는 repo.iris.tools 의 인증 정보이다.
                    docker.withRegistry("https://${env.IMAGE_REPO}", 'mobigen-harbor') {
                        // build/Dockerfile 의 멀티스테이지 빌드를 사용하여 Go 바이너리 빌드 후 최종 이미지 생성
                        def image = docker.build("${fullImageName}:${gitHash}", "-f build/Dockerfile .")

                        // Git hash 태그로 푸시
                        image.push()

                        // Latest 태그 추가 및 푸시
                        image.push("latest")
                    }

                    echo "Docker 이미지 빌드 및 푸시 완료: ${fullImageName}:${gitHash}"
                }
            }
        }
    }

    post {
        success {
            echo '빌드가 성공적으로 완료되었습니다!'
        }
        failure {
            echo '빌드가 실패했습니다.'
        }
        always {
            echo '빌드 완료'
        }
    }
}
