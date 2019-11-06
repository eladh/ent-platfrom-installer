module controller

require (
	common v0.0.0
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.9 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/spotahome/kooper v0.6.0
	golang.org/x/oauth2 v0.0.0-20190614102709-0f29369cfe45 // indirect
	golang.org/x/time v0.0.0-20190513212739-9d24e82272b4 // indirect
	k8s.io/api v0.0.0-20190111032252-67edc246be36
	k8s.io/apimachinery v0.0.0-20190313115320-c9defaaddf6f
	k8s.io/client-go v10.0.0+incompatible
)

replace common => ../common
