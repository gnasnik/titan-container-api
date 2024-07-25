package api

import (
	"fmt"
	ctypes "github.com/Filecoin-Titan/titan-container/api/types"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/docker/go-units"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-container-api/config"
	"github.com/gnasnik/titan-container-api/core/errors"
	"github.com/gnasnik/titan-container-api/core/generated/model"
	"net/http"
	"strconv"
	"strings"
)

func GetProvidersHandler(c *gin.Context) {
	url := config.Cfg.ContainerManager.Addr
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)

	params := ctypes.GetProviderOption{
		State: []ctypes.ProviderState{ctypes.ProviderStateOnline, ctypes.ProviderStateOffline, ctypes.ProviderStateAbnormal},
		Page:  int(page),
		Size:  int(size),
	}

	providers, err := getProvidersJsonRPC(url, params)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	type result struct {
		ID      string `json:"id"`
		IP      string `json:"ip"`
		State   string `json:"state"`
		Host    string `json:"host"`
		CPU     string `json:"cpu"`
		Memory  string `json:"memory"`
		Storage string `json:"storage"`
		Region  string `json:"region"`
	}

	res := make([]result, 0)
	for _, provider := range providers {
		resource, err := getProviderStatisticJsonRPC(url, provider.ID)
		if err != nil {
			log.Errorf("get statistic %s: %v", provider.ID, err)
			continue
		}

		//location, err := geo.GetIpLocation(c.Request.Context(), provider.IP)
		//if err != nil {
		//	log.Errorf("get location: %v", err)
		//}

		//if location == nil {
		//	location = &model.Location{}
		//}

		location := &model.Location{}

		res = append(res, result{
			ID:      string(provider.ID),
			IP:      provider.IP,
			State:   ctypes.ProviderStateString(provider.State),
			Host:    provider.HostURI,
			CPU:     fmt.Sprintf("%.1f/%.1f", resource.CPUCores.Available, resource.CPUCores.MaxCPUCores),
			Memory:  fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Memory.Available)), units.BytesSize(float64(resource.Memory.MaxMemory))),
			Storage: fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Storage.Available)), units.BytesSize(float64(resource.Storage.MaxStorage))),
			Region:  fmt.Sprintf("%s %s", location.Country, location.City),
		})
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"providers": res,
	}))
}

func GetDeploymentsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)

	params := ctypes.GetDeploymentOption{
		Owner: username,
		Page:  int(page),
		Size:  int(size),
	}

	resp, err := getDeploymentsJsonRPC(url, params)
	if err != nil {
		log.Errorf("get deployments: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	//type result struct {
	//	ID          string  `json:"id"`
	//	Name        string  `json:"name"`
	//	Image       string  `json:"image"`
	//	State       string  `json:"state"`
	//	Total       int     `json:"total"`
	//	Ready       int     `json:"ready"`
	//	Available   int     `json:"available"`
	//	CPU         float64 `json:"cpu"`
	//	GPU         float64 `json:"gpu"`
	//	Memory      string  `json:"memory"`
	//	Storage     string  `json:"storage"`
	//	Provider    string  `json:"provider"`
	//	Port        string  `json:"port"`
	//	CreatedTime string  `json:"created_time"`
	//}
	//
	//out := make([]ctypes.GetDeploymentListResp, 0)
	//
	//for _, deployment := range resp.Deployments {
	//	for _, service := range deployment.Services {
	//		state := ctypes.DeploymentStateInActive
	//		if service.Status.TotalReplicas == service.Status.ReadyReplicas {
	//			state = ctypes.DeploymentStateActive
	//		}
	//
	//		var exposePorts []string
	//		for _, port := range service.Ports {
	//			exposePorts = append(exposePorts, fmt.Sprintf("%d->%d", port.Port, port.ExposePort))
	//		}
	//
	//		var storageSize int64
	//		for _, storage := range service.Storage {
	//			storageSize += storage.Quantity
	//		}
	//
	//		out = append(out, result{
	//			ID:          string(deployment.ID),
	//			Name:        deployment.Name,
	//			Image:       service.Image,
	//			State:       ctypes.DeploymentStateString(state),
	//			Total:       service.Status.TotalReplicas,
	//			Ready:       service.Status.ReadyReplicas,
	//			Available:   service.Status.AvailableReplicas,
	//			CPU:         service.CPU,
	//			Memory:      units.BytesSize(float64(service.Memory * units.MiB)),
	//			Storage:     units.BytesSize(float64(storageSize * units.MiB)),
	//			Provider:    string(deployment.ProviderID),
	//			Port:        strings.Join(exposePorts, " "),
	//			CreatedTime: deployment.CreatedAt.Format(time.DateTime),
	//		})
	//	}
	//}

	c.JSON(http.StatusOK, respJSON(resp))
}

func GetDeploymentManifestHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	params := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	resp, err := getDeploymentsJsonRPC(url, params)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"deployment": resp.Deployments[0],
	}))
}

func CreateDeploymentHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	var deployment ctypes.Deployment
	err := c.BindJSON(&deployment)
	if err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusBadRequest, respError(errors.ErrInvalidParams))
		return
	}

	deployment.Owner = username
	url := config.Cfg.ContainerManager.Addr
	err = createDeploymentsJsonRPC(url, deployment)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusOK, respError(errors.ErrInvalidDeploymentName))
			return
		}

		log.Errorf("create deployment: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func DeleteDeploymentHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	var deployment ctypes.Deployment
	if err := c.BindJSON(&deployment); err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusBadRequest, respError(errors.ErrInvalidParams))
		return
	}

	deployment.Owner = username
	url := config.Cfg.ContainerManager.Addr
	err := deleteDeploymentsJsonRPC(url, deployment)
	if err != nil {
		log.Errorf("delete deployment: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func UpdateDeploymentHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	var deployment ctypes.Deployment
	if err := c.BindJSON(&deployment); err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusBadRequest, respError(errors.ErrInvalidParams))
		return
	}

	deployment.Owner = username
	url := config.Cfg.ContainerManager.Addr
	err := updateDeploymentsJsonRPC(url, deployment)
	if err != nil {
		log.Errorf("update deployment: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func GetDeploymentLogsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	params := ctypes.Deployment{
		Owner: username,
		ID:    ctypes.DeploymentID(deploymentId),
	}

	logs := make([]*ctypes.ServiceLog, 0)

	events, err := getDeploymentEventsJsonRPC(url, params)
	if err != nil {
		log.Errorf("get events: %v", err)
		//c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		//return
	}

	for _, event := range events {
		l := &ctypes.ServiceLog{
			ServiceName: event.ServiceName,
		}
		for _, e := range event.Events {
			l.Logs = append(l.Logs, ctypes.Log(e))
		}
		logs = append(logs, l)
	}

	slogs, err := getDeploymentLogsJsonRPC(url, params)
	if err != nil {
		log.Errorf("get logs: %v", err)
		//c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		//return
	}

	logs = append(logs, slogs...)

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"logs": logs,
	}))
}

func GetDeploymentEventsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	params := ctypes.Deployment{
		Owner: username,
		ID:    ctypes.DeploymentID(deploymentId),
	}

	logs, err := getDeploymentEventsJsonRPC(url, params)
	if err != nil {
		log.Errorf("get events: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"events": logs,
	}))
}

func GetDeploymentDomainHandler(c *gin.Context) {
	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	out := make([]*ctypes.DeploymentDomain, 0)
	domains, err := getDeploymentDomainJsonRPC(url, ctypes.DeploymentID(deploymentId))
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Errorf("get domains: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	out = append(out, domains...)
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"domains": out,
	}))
}

func AddDeploymentDomainHandler(c *gin.Context) {
	url := config.Cfg.ContainerManager.Addr
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	type domainReq struct {
		ID          ctypes.DeploymentID
		Hostname    string
		PrivateKey  string
		Certificate string
	}

	var params domainReq
	if err := c.BindJSON(&params); err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInvalidParams))
		return
	}

	dparam := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: params.ID,
	}

	resp, err := getDeploymentsJsonRPC(url, dparam)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	if len(resp.Deployments) == 0 {
		c.JSON(http.StatusOK, respJSON(errors.ErrNotFound))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrPermissionNotAllowed))
		return
	}

	cert := &ctypes.Certificate{
		Hostname:    params.Hostname,
		PrivateKey:  []byte(params.PrivateKey),
		Certificate: []byte(params.Certificate),
	}

	err = addDeploymentDomainJsonRPC(url, params.ID, cert)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Errorf("add domains: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func DeleteDeploymentDomainHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")
	host := c.Query("host")

	dparam := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	resp, err := getDeploymentsJsonRPC(url, dparam)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	if len(resp.Deployments) == 0 {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrPermissionNotAllowed))
		return
	}

	err = deleteDeploymentDomainJsonRPC(url, ctypes.DeploymentID(deploymentId), host)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Errorf("delete domains: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func GetDeploymentShellHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	params := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	resp, err := getDeploymentsJsonRPC(url, params)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	if len(resp.Deployments) == 0 {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrNotFound))
		return
	}

	shell, err := getDeploymentShellJsonRPC(url, ctypes.DeploymentID(deploymentId))
	if err != nil && !strings.Contains(err.Error(), "not found") {
		log.Errorf("get shell: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"shell": shell,
	}))
}

func GetIngressHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	dparam := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	resp, err := getDeploymentsJsonRPC(url, dparam)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrPermissionNotAllowed))
		return
	}

	ingress, err := getIngressJsonRPC(url, ctypes.DeploymentID(deploymentId))
	if err != nil {
		log.Errorf("get events: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"ingress": ingress,
	}))
}

func UpdateIngressHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	var ingress ctypes.Ingress
	if err := c.Bind(&ingress); err != nil {
		log.Errorf("err", err)
		c.JSON(http.StatusOK, respError(errors.ErrInvalidParams))
		return
	}

	dparam := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	resp, err := getDeploymentsJsonRPC(url, dparam)
	if err != nil {
		log.Errorf("get providers: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrPermissionNotAllowed))
		return
	}

	err = updateIngressJsonRPC(url, ctypes.DeploymentID(deploymentId), ingress.Annotations)
	if err != nil {
		log.Errorf("get ingress: %v", err)
		c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}
