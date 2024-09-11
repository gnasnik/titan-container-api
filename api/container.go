package api

import (
	"context"
	"fmt"
	ctypes "github.com/Filecoin-Titan/titan/api/types"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/docker/go-units"
	"github.com/gin-gonic/gin"
	"github.com/gnasnik/titan-container-api/core/dao"
	"github.com/gnasnik/titan-container-api/core/errors"
	"github.com/gnasnik/titan-container-api/core/generated/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAreasHandler(c *gin.Context) {
	areaIds := []string{"ALL"}

	for _, scheduler := range GlobalServer.GetSchedulers() {
		areaIds = append(areaIds, scheduler.AreaId)
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"area_ids": areaIds,
	}))
}

func GetProvidersHandler(c *gin.Context) {
	areaId := c.Query("area_id")
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)
	providerId := c.Query("provider_id")

	total, providers, err := dao.GetProvidersWithResource(c.Request.Context(), areaId, model.QueryOption{Page: int(page), Size: int(size), ID: providerId})
	if err != nil {
		log.Errorf("GetProvidersWithResource: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"total":     total,
		"providers": providers,
	}))
	return

	//option := &ctypes.GetProviderOption{
	//	State: []ctypes.ProviderState{ctypes.ProviderStateOnline},
	//	Page:  int(page),
	//	Size:  int(size),
	//}

	//scheduler, err := GetSchedulerByAreaId(areaId)
	//if err != nil {
	//	c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
	//	return
	//}
	//
	//response, err := scheduler.Api.GetProviderList(c.Request.Context(), option)
	//if err != nil {
	//	c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
	//	return
	//}
	//
	//out, err := syncQueryResource(context.Background(), scheduler, response)
	//if err != nil {
	//	c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
	//	return
	//}
	//
	//total := len(response)
	//if int64(total) > size {
	//	total = 50
	//}
	//
	//sort.Slice(out, func(i, j int) bool {
	//	return out[i].ID > out[j].ID
	//})
	//
	//c.JSON(http.StatusOK, respJSON(JsonObject{
	//	"total":     total,
	//	"providers": out,
	//}))
	//return

	//res := make([]result, 0)
	//
	//for _, provider := range response {
	//	//if provider.ResourcesStatistics != nil {
	//	//	resource = provider.ResourcesStatistics
	//	//}
	//
	//	//r := result{
	//	//	ID:         string(provider.ID),
	//	//	IP:         provider.IP,
	//	//	State:      ctypes.ProviderStateString(provider.State),
	//	//	RemoteAddr: provider.RemoteAddr,
	//	//}
	//
	//	//go func(res *result) {
	//	//	resource, err := scheduler.Api.GetStatistics(context.Background(), provider.ID)
	//	//	if err != nil {
	//	//		log.Errorf("get statistics: %v", err)
	//	//	}
	//	//	res.CPU = fmt.Sprintf("%.1f/%.1f", resource.CPUCores.Available, resource.CPUCores.MaxCPUCores)
	//	//	res.Memory = fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Memory.Available)), units.BytesSize(float64(resource.Memory.MaxMemory)))
	//	//	res.Storage = fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Storage.Available)), units.BytesSize(float64(resource.Storage.MaxStorage)))
	//	//
	//	//	ch <- res
	//	//}(&r)
	//
	//	resource, err := scheduler.Api.GetStatistics(context.Background(), provider.ID)
	//	if err != nil {
	//		log.Errorf("get statistics: %v", err)
	//	}
	//
	//	if resource == nil {
	//		resource = &ctypes.ResourcesStatistics{}
	//	}
	//
	//	location := &model.Location{}
	//	res = append(res, result{
	//		ID:         string(provider.ID),
	//		IP:         provider.IP,
	//		State:      ctypes.ProviderStateString(provider.State),
	//		RemoteAddr: provider.RemoteAddr,
	//		CPU:        fmt.Sprintf("%.1f/%.1f", resource.CPUCores.Available, resource.CPUCores.MaxCPUCores),
	//		Memory:     fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Memory.Available)), units.BytesSize(float64(resource.Memory.MaxMemory))),
	//		Storage:    fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Storage.Available)), units.BytesSize(float64(resource.Storage.MaxStorage))),
	//		Region:     fmt.Sprintf("%s %s", location.Country, location.City),
	//	})
	//}
	////
	////var out []*result
	//
	//c.JSON(http.StatusOK, respJSON(JsonObject{
	//	"total":     30,
	//	"providers": res,
	//}))
}

func syncQueryResource(ctx context.Context, scheduler *Scheduler, providers []*ctypes.Provider) ([]*model.ProviderWithResource, error) {
	var ch = make(chan *model.ProviderWithResource, 1)
	var out []*model.ProviderWithResource

	for _, provider := range providers {
		r := model.ProviderWithResource{
			ID:         provider.ID,
			Ip:         provider.IP,
			State:      int32(provider.State),
			RemoteAddr: provider.RemoteAddr,
			AreaID:     scheduler.AreaId,
			CreatedAt:  time.Now(),
		}

		go func(res *model.ProviderWithResource) {
			resource, err := scheduler.Api.GetStatistics(context.Background(), res.ID)
			if err != nil {
				log.Errorf("get statistics: %v", err)
			}

			if resource == nil {
				resource = &ctypes.ResourcesStatistics{}
			}

			res.Cpu = fmt.Sprintf("%.1f/%.1f", resource.CPUCores.Available, resource.CPUCores.MaxCPUCores)
			res.Memory = fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Memory.Available)), units.BytesSize(float64(resource.Memory.MaxMemory)))
			res.Storage = fmt.Sprintf("%s/%s", units.BytesSize(float64(resource.Storage.Available)), units.BytesSize(float64(resource.Storage.MaxStorage)))

			ch <- res
		}(&r)
	}

	for {
		select {
		case res := <-ch:
			out = append(out, res)

			if len(out) == len(providers) {
				return out, nil
			}
		}
	}

}

func GetDeploymentsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	areaId := c.Query("area_id")
	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(c.Query("size"), 10, 64)

	option := &ctypes.GetDeploymentOption{
		Owner: username,
		Page:  int(page),
		Size:  int(size),
	}

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), option)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(resp))
}

func GetDeploymentManifestHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	areaId := c.Query("area_id")
	deploymentId := c.Query("id")

	option := &ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), option)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	//resp, err := getDeploymentsJsonRPC(url, params)
	//if err != nil {
	//	log.Errorf("get providers: %v", err)
	//	c.JSON(http.StatusOK, respError(errors.ErrInternalServer))
	//	return
	//}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"deployment": resp.Deployments[0],
	}))
}

func CreateDeploymentHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	type createParams struct {
		ctypes.Deployment
		AreaId string
	}

	var deployment createParams
	err := c.BindJSON(&deployment)
	if err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusBadRequest, respError(errors.ErrInvalidParams))
		return
	}

	deployment.Owner = username
	//areaId := c.Query("area_id")

	scheduler, err := GetSchedulerByAreaId(deployment.AreaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	err = scheduler.Api.CreateDeployment(c.Request.Context(), &deployment.Deployment)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			c.JSON(http.StatusOK, respError(errors.ErrInvalidDeploymentName))
			return
		}

		log.Errorf("create deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func DeleteDeploymentHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	type deleteReq struct {
		Id     string `json:"id"`
		AreaId string `json:"area_id"`
	}

	var req deleteReq
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusBadRequest, respError(errors.ErrInvalidParams))
		return
	}

	scheduler, err := GetSchedulerByAreaId(req.AreaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	deployment := &ctypes.Deployment{
		ID:    ctypes.DeploymentID(req.Id),
		Owner: username,
	}
	err = scheduler.Api.CloseDeployment(c.Request.Context(), deployment, true)
	if err != nil {
		log.Errorf("delete deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func UpdateDeploymentHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	type updateReq struct {
		ctypes.Deployment
		AreaId string
	}

	var req updateReq
	if err := c.BindJSON(&req); err != nil {
		log.Errorf("%v", err)
		c.JSON(http.StatusBadRequest, respError(errors.ErrInvalidParams))
		return
	}

	req.Owner = username

	fmt.Printf("%+v\n", req)

	scheduler, err := GetSchedulerByAreaId(req.AreaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	err = scheduler.Api.UpdateDeployment(c.Request.Context(), &req.Deployment)
	if err != nil {
		log.Errorf("update deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func GetDeploymentLogsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	areaId := c.Query("area_id")

	//url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")

	params := ctypes.Deployment{
		Owner: username,
		ID:    ctypes.DeploymentID(deploymentId),
	}

	logs := make([]*ctypes.ServiceLog, 0)

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	events, err := scheduler.Api.GetEvents(c.Request.Context(), &params)
	if err != nil {
		log.Errorf("get event: %v", err)
		//c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		//return
		events = []*ctypes.ServiceEvent{{Events: []ctypes.Event{ctypes.Event(err.Error())}}}
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

	slogs, err := scheduler.Api.GetLogs(c.Request.Context(), &params)
	if err != nil {
		log.Errorf("get logs: %v", err)
		//c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		//return
		slogs = []*ctypes.ServiceLog{{Logs: []ctypes.Log{ctypes.Log(err.Error())}}}
	}

	logs = append(logs, slogs...)

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"logs": logs,
	}))
}

func GetDeploymentEventsHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	areaId := c.Query("area_id")
	deploymentId := c.Query("id")

	params := ctypes.Deployment{
		Owner: username,
		ID:    ctypes.DeploymentID(deploymentId),
	}

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	logs, err := scheduler.Api.GetLogs(c.Request.Context(), &params)
	if err != nil {
		log.Errorf("get logs: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"events": logs,
	}))
}

func GetDeploymentDomainHandler(c *gin.Context) {
	//url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")
	areaId := c.Query("area_id")

	out := make([]*ctypes.DeploymentDomain, 0)

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	domains, err := scheduler.Api.GetDeploymentDomains(c.Request.Context(), ctypes.DeploymentID(deploymentId))
	if err != nil {
		log.Errorf("get domain: %v", err)
		//c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		//return
	}

	out = append(out, domains...)
	c.JSON(http.StatusOK, respJSON(JsonObject{
		"domains": out,
	}))
}

func AddDeploymentDomainHandler(c *gin.Context) {
	//url := config.Cfg.ContainerManager.Addr
	areaId := c.Query("area_id")
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

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), &dparam)
	if err != nil {
		log.Errorf("get deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
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

	err = scheduler.Api.AddDeploymentDomain(c.Request.Context(), params.ID, cert)
	if err != nil {
		log.Errorf("add domain: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func DeleteDeploymentDomainHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	areaId := c.Query("area_id")
	deploymentId := c.Query("id")
	host := c.Query("host")

	dparam := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), &dparam)
	if err != nil {
		log.Errorf("get deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
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

	//}
	err = scheduler.Api.DeleteDeploymentDomain(c.Request.Context(), ctypes.DeploymentID(deploymentId), host)
	if err != nil {
		log.Errorf("delete domain: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}

func GetDeploymentShellHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)
	areaId := c.Query("area_id")
	deploymentId := c.Query("id")

	params := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), &params)
	if err != nil {
		log.Errorf("get deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
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

	leaseEndpoint, err := scheduler.Api.GetLeaseShellEndpoint(c.Request.Context(), ctypes.DeploymentID(deploymentId))
	if err != nil {
		log.Errorf("get deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"endpoint": leaseEndpoint,
	}))
}

func GetIngressHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	areaId := c.Query("area_id")
	deploymentId := c.Query("id")

	dparam := ctypes.GetDeploymentOption{
		Owner:        username,
		DeploymentID: ctypes.DeploymentID(deploymentId),
	}

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), &dparam)
	if err != nil {
		log.Errorf("get deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrPermissionNotAllowed))
		return
	}

	ingress, err := scheduler.Api.GetIngress(c.Request.Context(), ctypes.DeploymentID(deploymentId))
	if err != nil {
		log.Errorf("get ingress: %v", err)
		//c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		//return
	}

	c.JSON(http.StatusOK, respJSON(JsonObject{
		"ingress": ingress,
	}))
}

func UpdateIngressHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims[identityKey].(string)

	//url := config.Cfg.ContainerManager.Addr
	deploymentId := c.Query("id")
	areaId := c.Query("area_id")

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

	scheduler, err := GetSchedulerByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	resp, err := scheduler.Api.GetDeploymentList(c.Request.Context(), &dparam)
	if err != nil {
		log.Errorf("get deployment: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	if resp.Deployments[0].Owner != username {
		c.JSON(http.StatusOK, respError(errors.ErrPermissionNotAllowed))
		return
	}

	err = scheduler.Api.UpdateIngress(c.Request.Context(), ctypes.DeploymentID(deploymentId), ingress.Annotations)
	if err != nil {
		log.Errorf("update ingress: %v", err)
		c.JSON(http.StatusOK, respErrorWrapMessage(errors.ErrInternalServer, err.Error()))
		return
	}

	c.JSON(http.StatusOK, respJSON(nil))
}
