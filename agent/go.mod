module agent

require (
	common v0.0.0
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/hashicorp/go-version v1.1.0
	github.com/kris-nova/logger v0.0.0-20181127235838-fd0d87064b06
	github.com/thoas/go-funk v0.0.0-20181020164546-fbae87fb5b5c
	github.com/urfave/cli v1.20.0
	golang.org/x/oauth2 v0.0.0-20190624143730-0f29369cfe45 // indirect
	golang.org/x/time v0.0.0-20190513212739-9d24e82272b4 // indirect
	google.golang.org/genproto v0.0.0-20190819205937-24fa4b261c55 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace common => ../common
