package kubernetes

import (
	"bytes"
	"common/utils"
	"fmt"
	"github.com/kris-nova/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "k8s.io/api/core/v1"
	ext_v1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func GetClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		err1 := err
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			err = fmt.Errorf("InClusterConfig as well as BuildConfigFromFlags Failed. Error in InClusterConfig: %+v\nError in BuildConfigFromFlags: %+v", err1, err)
			return nil, err
		}
	}

	return config, nil
}

func GetClientsetFromConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		err = fmt.Errorf("failed creating clientset. Error: %+v", err)
		return nil, err
	}

	return clientset, nil
}

func GetClientset() (*kubernetes.Clientset, error) {
	config, err := GetClientConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		err = fmt.Errorf("failed creating clientset. Error: %+v", err)
		return nil, err
	}

	return clientset, nil
}

func ExecOnPod(command, podName, containerName, namespace string, stdin io.Reader) (string, string, error) {
	config, err := GetClientConfig()
	if err != nil {
		return "", "", err
	}

	clientset, err := GetClientsetFromConfig(config)
	if err != nil {
		return "", "", err
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&PodExecOptions{
		Command:   strings.Fields(command),
		Container: containerName,
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	logger.Always(stdout.String())
	logger.Always(stderr.String())

	return stdout.String(), stderr.String(), nil
}

func UploadFileToPod(file string, targetPath string, podName string, containerName string, namespace string) (string, string, error) {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		logger.Critical("get resource error", err)
	}

	return UploadContentToPod(content, targetPath, podName, containerName, namespace)
}

func UploadContentToPod(content []byte, targetPath string, podName string, containerName string, namespace string) (string, string, error) {
	stdin := bytes.NewReader(content)
	return ExecOnPod("cp /dev/stdin "+ targetPath, podName, containerName, namespace, stdin)
}

func DownloadFileFromPod(file string, podName string, containerName string, namespace string) (string, string, error) {
	return ExecOnPod("cat  "+file, podName, containerName, namespace, nil)
}

func DeleteSecret(secretName string, namespace string) error {
	clientset, err := GetClientset()
	if err != nil {
		return err
	}

	_ = clientset.CoreV1().Secrets(namespace).Delete(secretName, nil)

	return nil
}

func DeleteConfigMap(configmapName string, namespace string) error {
	clientset, err := GetClientset()
	if err != nil {
		return err
	}

	_ = clientset.CoreV1().ConfigMaps(namespace).Delete(configmapName, nil)

	return nil
}

func CreateConfigMap(configmapName string, fileName string,namespace string, content string) error {
	clientset, err := GetClientset()

	if err != nil {
		return err
	}

	var configMap = make(map[string]string)
	configMap[fileName] = content

	_, err = clientset.CoreV1().ConfigMaps(namespace).Create(&ConfigMap{
		Data: configMap,
		ObjectMeta: metav1.ObjectMeta{
			Name: configmapName,
		},
	})

	return nil
}

func GetSecret(name string ,attribute string) string {
	clientset, _ := GetClientset()
	res, _:= clientset.CoreV1().Secrets("default").Get(name ,metav1.GetOptions {})

	return string(res.Data[attribute])
}

func CreateSecret(secretName string, namespace string, publicKey string, privateKey string) error {
	clientset, err := GetClientset()

	if err != nil {
		return err
	}

	var secretMap = make(map[string]string)
	secretMap["tls.key"] = privateKey
	secretMap["tls.crt"] = publicKey

	_, err = clientset.CoreV1().Secrets(namespace).Create(&Secret{
		Type:       SecretType("kubernetes.io/tls"),
		StringData: secretMap,
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
	})

	return nil
}

func GetServiceExternalAddress(serviceName string) (string, error) {
	clientset, err := GetClientset()
	if err != nil {
		return "", nil
	}

	services, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})

	if err != nil {
		return "", nil
	}

	for _, service := range services.Items {
		if strings.Contains(service.Name, serviceName) {
			if service.Status.LoadBalancer.Ingress == nil {
				return "", nil
			}
			ingress := service.Status.LoadBalancer.Ingress[0]
			if ingress.IP != "" {
				return ingress.IP, nil
			} else {
				return ingress.Hostname, nil
			}
		}
	}

	return "", nil
}

func GetPod(name string) (Pod, error) {
	clientset, err := GetClientset()
	if err != nil {
		return Pod{}, err
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return Pod{}, err
	}
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, name) {
			return pod, nil
		}
	}

	return Pod{} , err
}

func GetNodes() (*NodeList, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, nil
	}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	return nodes, nil
}

func GetDeployments() (*ext_v1.DeploymentList, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, nil
	}

	deployments, err := clientset.ExtensionsV1beta1().Deployments("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	return deployments, nil
}


func GetDaemonset() (*ext_v1.DaemonSetList, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, nil
	}

	daemonsets, err := clientset.ExtensionsV1beta1().DaemonSets("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	return daemonsets, nil
}


func GetIngress() (*ext_v1.IngressList, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, nil
	}

	ingress, err := clientset.ExtensionsV1beta1().Ingresses("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	return ingress, nil
}


func GetServices() (*ServiceList, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, nil
	}

	services, err := clientset.CoreV1().Services("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	return services, nil
}


func GetPods() (*PodList, error) {
	clientset, err := GetClientset()
	if err != nil {
		return nil, nil
	}

	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	return pods, nil
}

func CreateNamespace(name string) error {
	clientset, err := GetClientset()
	if err != nil {
		return nil
	}

	nsSpec := &Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	clientset.CoreV1()

	_, err = clientset.CoreV1().Namespaces().Create(nsSpec)

	return err
}


func WaitForPod(name string ,containers int) {
	logger.Always("wait for pod" , name)
	podReady := false
	for !podReady {
		pod, err := GetPod(name)
		if err == nil && pod.Name != "" &&  len(pod.Status.ContainerStatuses) == containers {
			podReady = true
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}


func GetServiceIp(serviceName string) string {
	//todo - timeout for service discovery
	address := ""
	for address == "" {
		addressCandidate, err := GetServiceExternalAddress(serviceName)
		if err == nil && addressCandidate != "" {
			address = addressCandidate
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	return address
}

func ExposePod(name string ,serviceType string ,namespace string) {
	utils.Shell("kubectl expose pod " + name + " --type=" + serviceType + " --name=" + name +
		" --namespace=" + namespace)
}