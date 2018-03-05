#!/usr/bin/env groovy

/*
 
    HERE BE DRAGONS

    This one seed will create all microservices in a list defined below

*/
//  list of apps
// devlop / master branches for each pipeline
def pipelineNames = ["dev", "test"]
def microservices = ["notification-microservice"]

//  Globals for across all the jobs
def pipelineNamespace = "ci-cd"
newLine = System.getProperty("line.separator")

def pipelineGeneratorVersion = "${JOB_NAME}.${BUILD_ID}"

def jenkinsGitCreds = "jenkins-git-creds"

//  Common functions repeated across the jobs
def buildWrappers(context) {
    context.ansiColorBuildWrapper {
        colorMapName('xterm')
    }
}

def notifySlack(context) {
    context.slackNotifier {
        notifyAborted(true)
        notifyBackToNormal(true)
        notifyFailure(true)
        notifyNotBuilt(true)
        notifyRegression(true)
        notifyRepeatedFailure(true)
        notifySuccess(true)
        notifyUnstable(true)
    }
}

def rotateLogs(context) {
    context.logRotator {
        daysToKeep(100)
        artifactNumToKeep(2)
    }
}

def coverageReport(context, appName) {
    // TODO - fix thiss for the golang app
    // context.cobertura('coverage.xml') {
    //     failNoReports(true)
    //     sourceEncoding('ASCII')
    //     // the following targets are added by default to check the method, line and conditional level coverage
    //     methodTarget(80, 40, 20)
    //     lineTarget(80, 40, 20)
    //     conditionalTarget(70, 40, 20)
    // }
    context.publishHtml {
        report('src/redhat/' + appName + '/reports') {
            reportName('HTML Code Coverage Report')
            allowMissing(false)
            alwaysLinkToLastBuild(false)
        }
    }
}
microservices.each {
    def appName = it
    def gitBaseUrl = 'https://gitlab.apps.xpertex.rht-labs.com/xpertex/' + appName + '.git'

    pipelineNames.each {
        def pipelineName = it
        def buildImageName = it + "-" + appName + "-build"
        def bakeImageName = it + "-" + appName + "-bake"
        def deployImageName = it + "-" + appName + "-deploy"
        def zapSecurityScan = it + "-" + appName + "-zap-scan"
        def arachniSecurityScan = it + "-" + appName + "-arachni-scan"
        def projectNamespace = "labs-" + it
        def appEndpoint = 'http://'+ appName + "-" + projectNamespace  + '.apps.xpertex.rht-labs.com/'
        def jobDescription = "THIS JOB WAS GENERATED BY THE JENKINS SEED JOB - ${pipelineGeneratorVersion}.  \n"  + it + " backend build job for the app."

        job(buildImageName) {
            description(jobDescription)
            label('golang-build-pod')

            rotateLogs(delegate)

            wrappers {
                buildWrappers(delegate)

                preScmSteps {
                    // Fails the build when one of the steps fails.
                    failOnError(true)
                    // Adds build steps to be run before SCM checkout.
                    steps {
                        //  TODO - add git creds here
                        shell('git config --global http.sslVerify false' + newLine +
                                'git config --global user.name jenkins' + newLine +
                                'git config --global user.email jenkins@cc.net')
                    }
                }
            }
            scm {
                git {
                    remote {
                        name('origin')
                        url(gitBaseUrl)
                        credentials(jenkinsGitCreds)
                    }
                    if (pipelineName.contains('test')){
                        branch('master')
                    }
                    else {
                        branch('develop')
                    }
                }
            }
            if (pipelineName.contains('dev')){
                triggers {
                    cron('H/60 H/2 * * *')
                    gitlabPush {
                        buildOnPushEvents()
                    }
                }
            }
            steps {
                steps {
                    shell('#!/bin/bash' + newLine +
                            'set -o xtrace' + newLine +
                            'export GOPATH=${WORKSPACE}' + newLine +
                            'mkdir -p $GOPATH/src/redhat/' + appName + newLine +
                            'cp -r * $GOPATH/src/redhat/' + appName + newLine +
                            'set -e' + newLine +
                            'cd $GOPATH/src/redhat/' + appName + newLine +
                            'NAME=' + appName + newLine +
                            'go get github.com/onsi/ginkgo/ginkgo' + newLine +
                            'go get -v -t ./...' + newLine +
                            'go build -v' + newLine +
                            '$GOPATH/bin/ginkgo -r --cover -keepGoing' + newLine +
                            'echo "Converting ONE report to html - TODO them all"' + newLine +
                            'mkdir -p reports' + newLine +
                            'go tool cover -html=api/api.coverprofile  -o reports/index.html' + newLine +
                            'mv ' + appName + ' build' + newLine +
                            'mkdir package-contents' + newLine +
                            'mv Dockerfile build package-contents' + newLine +
                            'cp ' + projectNamespace + '.toml package-contents/config.toml' + newLine +
                            'zip -r ' + appName + '.zip package-contents')
                }
            }
            publishers {
                // nexus upload
                postBuildScripts {
                    steps {
                        shell('curl -v -F r=releases \\' + newLine +
                                    '-F hasPom=false \\' + newLine +
                                    '-F e=zip \\' + newLine +
                                    '-F g=com.example.go \\' + newLine +
                                    '-F a=' + appName + ' \\' + newLine +
                                    '-F v=0.0.1-${JOB_NAME}.${BUILD_NUMBER} \\' + newLine +
                                    '-F p=zip \\' + newLine +
                                    '-F file=@${WORKSPACE}/src/redhat/' + appName + '/' + appName + '.zip \\' + newLine +
                                    '-u admin:admin123 http://nexus-v2.ci-cd.svc.cluster.local:8081/nexus/service/local/artifact/maven/content')
                    }
                }

                archiveArtifacts('**')

                coverageReport(delegate, appName)

                xUnitPublisher {
                    tools {
                        jUnitType {
                            pattern('**/test-report.xml')
                            skipNoTestFiles(false)
                            failIfNotNew(true)
                            deleteOutputFiles(true)
                            stopProcessingIfError(true)
                        }
                    }
                    thresholds {
                        failedThreshold {
                            failureThreshold('0')
                            unstableThreshold('')
                            unstableNewThreshold('')
                            failureNewThreshold('')
                        }
                    }

                    thresholdMode(0)
                    testTimeMargin('3000')
                }
                // git publisher
                git {
                    tag("origin", "\${JOB_NAME}.\${BUILD_NUMBER}") {
                        create(true)
                        message("Automated commit by jenkins from \${JOB_NAME}.\${BUILD_NUMBER}")
                    }
                }

                downstreamParameterized {
                    trigger(bakeImageName) {
                        condition('UNSTABLE_OR_BETTER')
                        parameters {
                            predefinedBuildParameters{
                                properties("BUILD_TAG=\${JOB_NAME}.\${BUILD_NUMBER}")
                                textParamValueOnNewLine(true)
                            }
                        }
                    }
                }
                
                downstreamParameterized {
                    trigger("hue-fail") {
                        condition("FAILED")
                        triggerWithNoParameters()
                    }
                }

                notifySlack(delegate)
            }
        }


        job(bakeImageName) {
            description(jobDescription)
            parameters{
                string{
                    name("BUILD_TAG")
                    defaultValue("my-app-build.1234")
                    description("The BUILD_TAG is the \${JOB_NAME}.\${BUILD_NUMBER} of the successful build to be promoted.")
                }
            }
            rotateLogs(delegate)

            wrappers {
                buildWrappers(delegate)
            }
            steps {
                steps {
                    shell('#!/bin/bash' + newLine +
                            'set -o xtrace' + newLine +
                            '# WIPE PREVIOUS BINARY' + newLine +
                            'rm -rf *.zip package-contents' + newLine +
                            '# GET BINARY - DIRTY GET BINARY HACK' + newLine +
                            'curl -v -f http://admin:admin123@nexus-v2.ci-cd.svc.cluster.local:8081/nexus/service/local/repositories/releases/content/com/example/go/' + appName + '/0.0.1-${BUILD_TAG}/' + appName + '-0.0.1-${BUILD_TAG}.zip -o ' + appName + '.zip' + newLine +
                            'unzip ' + appName + newLine +
                            'oc project ci-cd'  + newLine +
                            '# DO OC BUILD STUFF WITH BINARY NOW' + newLine +
                            'NAME=' + appName  + newLine +
                            'oc patch bc ${NAME} -p "spec:' + newLine +
                            '   nodeSelector:' + newLine +
                            '   output:' + newLine +
                            '     to:' + newLine +
                            '       kind: ImageStreamTag' + newLine +
                            '       name: \'${NAME}:${JOB_NAME}.${BUILD_NUMBER}\'"' + newLine +
                            'oc start-build ${NAME} --from-dir=package-contents/ --follow')
                }
            }
            publishers {
                downstreamParameterized {
                    trigger(deployImageName) {
                        condition('SUCCESS')
                        parameters {
                            predefinedBuildParameters{
                                properties("BUILD_TAG=\${JOB_NAME}.\${BUILD_NUMBER}")
                                textParamValueOnNewLine(true)
                            }
                        }
                    }
                }
                downstreamParameterized {
                    trigger("hue-fail") {
                        condition("FAILED")
                        triggerWithNoParameters()
                    }
                }
                notifySlack(delegate)
            }
        }

        job(deployImageName) {
            description(jobDescription)
            parameters {
                string{
                    name("BUILD_TAG")
                    defaultValue("my-app-build.1234")
                    description("The BUILD_TAG is the \${JOB_NAME}.\${BUILD_NUMBER} of the successful build to be promoted.")
                }
            }
            rotateLogs(delegate)

            wrappers {
                buildWrappers(delegate)

            }
            steps {
                steps {
                    shell('#!/bin/bash' + newLine +
                            'set -o xtrace' + newLine +
                            'PIPELINES_NAMESPACE=' + pipelineNamespace  + newLine +
                            'NAMESPACE=' + projectNamespace  + newLine +
                            'NAME=' + appName  + newLine +
                            'oc tag ${PIPELINES_NAMESPACE}/${NAME}:${BUILD_TAG} ${NAMESPACE}/${NAME}:${BUILD_TAG}' + newLine +
                            'oc project ${NAMESPACE}' + newLine +
                            'oc patch dc ${NAME} -p "spec:' + newLine +
                            '  template:' + newLine +
                            '    spec:' + newLine +
                            '      containers:' + newLine +
                            '        - name: ${NAME}' + newLine +
                            '          image: \'docker-registry.default.svc:5000/${NAMESPACE}/${NAME}:${BUILD_TAG}\'' + newLine +
                            '          env:' + newLine +
                            '            - name: NODE_ENV' + newLine +
                            '              value: \'production\'"' + newLine +
                            'oc rollout latest dc/${NAME}')
                }
                openShiftDeploymentVerifier {
                    apiURL('')
                    depCfg(appName)
                    namespace(projectNamespace)
                    // This optional field's value represents the number expected running pods for the deployment for the DeploymentConfig specified.
                    replicaCount('1')
                    authToken('')
                    verbose('yes')
                    // This flag is the toggle for turning on or off the verification that the specified replica count for the deployment has been reached.
                    verifyReplicaCount('yes')
                    waitTime('')
                    waitUnit('sec')
                }
            }
            publishers {
                downstreamParameterized {
                    trigger(zapSecurityScan) {
                        condition('UNSTABLE_OR_BETTER')
                        triggerWithNoParameters()
                    }
                }
                downstreamParameterized {
                    trigger(arachniSecurityScan) {
                        condition('UNSTABLE_OR_BETTER')
                        triggerWithNoParameters()
                    }
                }
                downstreamParameterized {
                    trigger("hue-fail") {
                        condition("FAILED")
                        triggerWithNoParameters()
                    }
                }
                notifySlack(delegate)
            }
        }

        job(zapSecurityScan) {
            description(jobDescription)
            label('zap-build-pod')

            rotateLogs(delegate)

            wrappers {
                buildWrappers(delegate)
            }
            steps {
                shell {
                    command('#!/bin/bash' + newLine +
                            'export URL_TO_TEST=' + appEndpoint + newLine +
                            '/zap/zap-baseline.py -r index.html -t ${URL_TO_TEST}' + newLine +
                            'RC=$?' + newLine +
                            'echo "Zap exit code :: ${RC}"' + newLine +
                            'exit ${RC}')
                    unstableReturn(2)
                }
            }
            publishers {
                publishHtml {
                    report('/zap/wrk') {
                        reportName('Zap Security Report')
                        allowMissing(false)
                        alwaysLinkToLastBuild(false)
                    }
                }

                notifySlack(delegate)
            }
        }

        job(arachniSecurityScan) {
            description(jobDescription)
            label('arachni-build-pod')

            rotateLogs(delegate)

            wrappers {
                buildWrappers(delegate)
            }
            
            steps {
                shell {
                    command('#!/bin/bash' + newLine +
                            'set -o xtrace' + newLine +
                            'export URL_TO_TEST='+ appEndpoint + newLine +
                            '/arachni/bin/arachni ${URL_TO_TEST} --report-save-path=arachni-report.afr' + newLine +
                            '/arachni/bin/arachni_reporter arachni-report.afr --reporter=xunit:outfile=report.xml --reporter=html:outfile=web-report.zip' + newLine +
                            'unzip web-report.zip -d arachni-web-report')
                    unstableReturn(1)
                }
            }
            publishers {
                publishHtml {
                    report('arachni-web-report') {
                        reportName('Arachni Report')
                        allowMissing(false)
                        alwaysLinkToLastBuild(false)
                    }
                }
                xUnitPublisher {
                    tools {
                        jUnitType {
                            pattern('report.xml')
                            skipNoTestFiles(false)
                            failIfNotNew(true)
                            deleteOutputFiles(true)
                            stopProcessingIfError(true)
                        }
                    }
                    thresholds {
                        failedThreshold {
                            failureThreshold('0')
                            unstableThreshold('5')
                            unstableNewThreshold('5')
                            failureNewThreshold('0')
                        }
                    }

                    thresholdMode(0)
                    testTimeMargin('3000')
                }

                notifySlack(delegate)
            }
        }
        buildPipelineView(pipelineName + '-' + appName + "-pipeline") {
            filterBuildQueue()
            filterExecutors()
            title(pipelineName + ' ' + appName + " CI Pipeline")
            displayedBuilds(10)
            selectedJob(buildImageName)
            alwaysAllowManualTrigger()
            refreshFrequency(5)
        }

        buildMonitorView(appName +'-monitor') {
            description('All build jobs for the microservices')
            filterBuildQueue()
            filterExecutors()
            jobs {
                regex('.*' + appName + '.*')
            }
        }

    }
}
