package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

var K8sClient *unversioned.Client

type k8sConfig struct {
	Host     string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
	Insecure bool   `default:"false"`
}

type AppStatus struct {
	Namespace         string
	MinReplicas       int32
	MaxReplicas       int32
	TargetCPU         string
	LastScaleTime     string
	DesiredReplicas   int32
	CurrentReplicas   int32
	CurrentCPU        string
	UpToDateReplicas  int32
	AvailableReplicas int32
}

func main() {
	router := gin.Default()

	auth := router.Group("/", gin.BasicAuth(gin.Accounts{
		"leroy": os.Getenv("PASSWORD"),
	}))

	auth.GET("/", func(c *gin.Context) {
		m := make(map[string]AppStatus)

		// hpa
		hpalist, _ := K8sClient.HorizontalPodAutoscalers(api.NamespaceAll).List(api.ListOptions{})
		for _, h := range hpalist.Items {
			app := AppStatus{}
			app.Namespace = h.Namespace
			app.MinReplicas = *h.Spec.MinReplicas
			app.MaxReplicas = h.Spec.MaxReplicas
			if h.Spec.TargetCPUUtilizationPercentage != nil {
				app.TargetCPU = fmt.Sprint(*h.Spec.TargetCPUUtilizationPercentage, "%")
			}
			app.CurrentReplicas = h.Status.CurrentReplicas
			app.DesiredReplicas = h.Status.DesiredReplicas
			if h.Status.CurrentCPUUtilizationPercentage != nil {
				app.CurrentCPU = fmt.Sprint(*h.Status.CurrentCPUUtilizationPercentage, "%")
			}
			if h.Status.LastScaleTime != nil {
				app.LastScaleTime = h.Status.LastScaleTime.Local().String()
			}
			m[app.Namespace] = app
		}
		deplist, _ := K8sClient.Deployments(api.NamespaceAll).List(api.ListOptions{})
		for _, d := range deplist.Items {
			app, ok := m[d.Namespace]
			if !ok {
				app = AppStatus{}
				app.Namespace = d.Namespace
			}
			app.CurrentReplicas = (d.Status.Replicas - d.Status.UnavailableReplicas)
			app.UpToDateReplicas = d.Status.UpdatedReplicas
			app.AvailableReplicas = d.Status.AvailableReplicas
			m[d.Namespace] = app
		}

		apps := []AppStatus{}

		// creating regexp to filter namespaces...
		re := regexp.MustCompile(os.Getenv("NAMESPACE_FILTER"))
		for k, v := range m {
			if re.MatchString(k) {
				apps = append(apps, v)
			}
		}

		c.JSON(http.StatusOK, gin.H{"data": apps})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	router.Run(fmt.Sprint(":", port))
}

func init() {
	configenv := &k8sConfig{}
	err := envconfig.Process("k8s", configenv)
	if err != nil {
		log.Panicf("Failed to read k8s configuration from environment: %s", err.Error())
	}
	// K8s config
	config := &restclient.Config{
		Host:     configenv.Host,
		Username: configenv.Username,
		Password: configenv.Password,
		Insecure: configenv.Insecure,
	}
	// Creating k8s client
	K8sClient, err = unversioned.New(config)
	if err != nil {
		log.Panicf("Error trying to create a kubernetes client. Error: %s", err.Error())
	}
}
