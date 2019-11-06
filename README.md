## **Platform Installer:**
### **Modules:**

#### Agent
Installer - (Install k8s -> Install Helm -> Install E+)

#### Cli
Installer CLI

#### Common
Installer Common services/utils and structures

#### Controller
Installer K8s Agent - REST API + Web Application (React) - suppose to be deployed to target cluster


**remember to use the new agent installer docker image tag in CLI _--version_ value*

### **Install new E+ setup using the CLI:**

```
./platformCli i \
 --config=/Users/userName/Desktop/workspace/consulting/install-platform-go/cli/resources/setup-eplus-full.yaml \
 --secret=/Users/userName/Desktop/devops-consulting-1572a82131fc.json \
 --name=name \
 --instances=1 \
 --version=2.0.15 \
 --artifactory=/Users/userName/Desktop/workspace/consulting/install-platform-go/agent/resources/buckets/artifactory.json \
 --edge=/Users/userName/Desktop/workspace/consulting/install-platform-go/agent/resources/buckets/edge.json
 ```
 
 ### **Uninstall E+ setup using the CLI:**
 
 ```
 ./platformCli u \
  --config=/Users/userName/Desktop/workspace/consulting/install-platform-go/cli/resources/setup-eplus-full.yaml \
  --secret=/Users/userName/Desktop/devops-consulting-1572a82131fc.json \
  --name=name \
  --instances=1 \
  --version=2.0.15 \
  ```
  
  ### **List registered E+ setup ssh servers using the CLI:**
   ```
 
   ./platformCli l \
    --config=/Users/userName/Desktop/workspace/consulting/install-platform-go/cli/resources/setup-eplus-full.yaml \
    --secret=/Users/userName/Desktop/devops-consulting-1572a82131fc.json \
    --name=name \
    --instances=1 \
    --version=2.0.15
   ```

 
### **Setup Yaml :**

#### cluster

name - cluster name
domain - domain name

#### vendor
type - gcp/aws/azure \
region - vendor region (us-central1) \
zone - vendor zone (us-central1-a) \
project - vendor project (devops-consulting)  


#### GCP vendor specific
storage.identity = gcp storage identity \
storage.secret =  gcp storage secret

#### License buckets
art_license = artifactory E+ servers license bucket code \
edge_license= artifactory E+ edge server license bucket code 


#### Sites
 ```
sites:
  - name: USA
    description: US West coast site
    city:
      name: Sunnyvale
      country_code: US
      latitude: 37.36883
      longitude: -122.03635
  - name: England
    description: Europe West site
    city:
      name: London
      country_code: GB
      latitude: 51.5074
      longitude: 0.1278
 ```


#### Services
manage e+ helm charts versions

 ```
  versions:
    artifactory: 7.17.1
    distribution: 3.4.0
    jfmc: 1.1.5
    xray: 1.0.5
    sonar: 0.15.0
    jenkins: 2.164.2
 ```


Artifactory 
 ```
  artifactory:
    - name: artifactory
      site: USA
      auth_server: true
      repos:
        - name: gradle
          local: true
          remote: true
          virtual: true
          url: https://jcenter.bintray.com
          package_type: gradle
        - name: docker
          local: true
          remote: true
          virtual: true
          url: https://registry-1.docker.io/
          package_type: docker
        - name: docker-prod
          local: true
          package_type: docker
        - name: npm
          local: true
          remote: true
          virtual: true
          url: https://registry.npmjs.org
          package_type: npm
        - name: helm
          local: true
 ```
 
 

Edge 
 ```
   edges:
     - name: edge-london
       site: England
       auth_server: false
       repos:
         - name: docker-prod
           local: true
         - name: helm
           local: true
 
  ```
  
Distribution 
 ```

    distribution:
      name: distribution
      site: USA
 ```


Xray

 ```

 xray:
    - name: xray-server
      site: USA
      artifactory: artifactory
      policies:
        - name: securityPolicy
          type: security
          description: some description
          rules:
            - name: securityRule
              priority: 1
              criteria:
                min_severity: all severities
              actions:
                fail_build: false
                block_download:
                  unscanned: false
                  active: true
      watches:
        - general_data:
            name: vuln-prod
            description: This is a watch for security threats
            active: true
          project_resources:
            resources:
              - type: repository
                bin_mgr_id: artifactory
                name: docker-local
                filters:
                  - type: regex
                    value: ".*"
              - type: build
                name: docker-app-demo  
                bin_mgr_id: artifactory
                clickable: true
          assigned_policies:
            - name: securityPolicy
              type: security
 ```
 
 
 Tools
  ```
  tools:
   dev: true
   sonarqube: true
   glowroot: true
   jenkins:
     site: USA
     jobs:
       - name: npm-app-demo
         url: https://github.com/jfrog/consulting
         pipeline: jenkins/npm-app-demo/Jenkinsfile.groovy
         params:
           - name: dddd
             type: booleanParam
             default_value: false
             desc: uncheck to disable tests
       - name: npm-app-demo2
         url: https://github.com/jfrog/consulting
         pipeline: jenkins/npm-app-demo/Jenkinsfile.groovy
  ```
 
