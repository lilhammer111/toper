pipeline {
    agent any

    stages {
        stage('pull code') {
            steps {
                withCredentials([string(credentialsId: 'github-token-for-gin-test', variable: 'TOKEN')]) {
                    git branch: 'main', credentialsId: 'github-email-pwd', url: 'https://${TOKEN}@github.com/lilhammer111/to-persist.git'
                }

            }
        }
        stage('build project') {
            steps {
                sh '''echo "start building..."
                    echo "finish building!"'''
            }
        }
        stage('deploy project') {
            steps {
                sshPublisher(publishers: [sshPublisherDesc(configName: 'qiniuyun', transfers: [sshTransfer(cleanRemote: false, excludes: '', execCommand: 'echo "success!"', execTimeout: 120000, flatten: false, makeEmptyDirs: false, noDefaultExcludes: false, patternSeparator: '[, ]+', remoteDirectory: '/to-persist', remoteDirectorySDF: false, removePrefix: '', sourceFiles: '**')], usePromotionTimestamp: false, useWorkspaceInPromotion: false, verbose: false)])
            }
        }
    }
}
