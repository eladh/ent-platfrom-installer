package installers

import (
	"agent/utils"
	"bytes"
	"common/kubernetes"
	"common/structs"
	commonUtils "common/utils"
	"github.com/kris-nova/logger"
)

const productionValue string = "https://raw.githubusercontent.com/helm/charts/master/stable/postgresql/values-production.yaml"
const singleNodeValue string = "https://raw.githubusercontent.com/helm/charts/master/stable/postgresql/values.yaml"

// TODO - add this lines to startup command
//CREATE USER artifactory WITH PASSWORD 'password';
//CREATE DATABASE artifactory WITH OWNER=artifactory ENCODING='UTF8';
//GRANT ALL PRIVILEGES ON DATABASE artifactory TO artifactory;

func InstallPostgresSingleNode(name string ,password string ,setupInfo structs.SetupInfo) {
	config := setupInfo.TempDir + "postgresql-values.yaml"
	_, _ = commonUtils.DownloadFile(config, singleNodeValue)

	commonUtils.HelmInstall(name ,"stable/postgresq" ,"",
		[]string{"postgresqlPassword=" + password , "replication.password=" + password}, config);

	utils.WaitForService(name, 80, true)
}

func SetupArtifactoryDb(name string ,password string) {
	var cmd bytes.Buffer
	// see details - https://github.com/helm/charts/issues/9619
	cmd.WriteString("PGPASSWORD=\"password\" psql -c \"CREATE USER artifactory WITH PASSWORD 'password'\"")
	result ,_ , _ := kubernetes.ExecOnPod(cmd.String(),  name + "-postgresql-0", "", "default", nil)
	logger.Always(result)

	//CREATE DATABASE artifactory WITH OWNER=artifactory ENCODING='UTF8';
	//GRANT ALL PRIVILEGES ON DATABASE artifactory TO artifactory;
}