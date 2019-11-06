package utils

import (
	"common/kubernetes"
	"github.com/kris-nova/logger"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var servicesWaitingGroup sync.WaitGroup

func WaitForAllServices() {
	servicesWaitingGroup.Wait()
	logger.Always("finished waiting for install")
}

func WaitForService(name string, port int, background bool) {
	if background {
		servicesWaitingGroup.Add(1)
		go waitForServiceInternal(&servicesWaitingGroup, name, port)
	} else {
		waitForServiceInternal(nil, name, port)
	}
}

func waitForServiceInternal(wg *sync.WaitGroup, name string, port int) {
	if wg != nil {
		defer wg.Done()
	}

	responseOk := false
	serviceIp := kubernetes.GetServiceIp(name)
	for !responseOk {
		resp, err := http.Get("http://" + serviceIp + ":" + strconv.Itoa(port))
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 403) {
			responseOk = true
			logger.Always("Service " + name + " is ready")
		} else {
			time.Sleep(20 * time.Second)
		}
	}
}

func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}