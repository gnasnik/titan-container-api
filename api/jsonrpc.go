package api

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	ctypes "github.com/Filecoin-Titan/titan/api/types"
//	"github.com/gnasnik/titan-container-api/core/generated/model"
//	"golang.org/x/xerrors"
//	"io/ioutil"
//	"net/http"
//)
//
//func getProvidersJsonRPC(url string, opt ctypes.GetProviderOption) ([]*ctypes.Provider, error) {
//	params, err := json.Marshal([]interface{}{opt})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetProviderList",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var providers []*ctypes.Provider
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &providers)
//	if err != nil {
//		return nil, err
//	}
//
//	return providers, nil
//}
//
//func getProviderStatisticJsonRPC(url string, id string) (*ctypes.ResourcesStatistics, error) {
//	params, err := json.Marshal([]interface{}{id})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetStatistics",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var statistic ctypes.ResourcesStatistics
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &statistic)
//	if err != nil {
//		return nil, err
//	}
//
//	return &statistic, nil
//}
//
//func getDeploymentsJsonRPC(url string, opt ctypes.GetDeploymentOption) (*ctypes.GetDeploymentListResp, error) {
//	params, err := json.Marshal([]interface{}{opt})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetDeploymentList",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var out ctypes.GetDeploymentListResp
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &out)
//	if err != nil {
//		return nil, err
//	}
//
//	return &out, nil
//}
//
//func createDeploymentsJsonRPC(url string, deployment ctypes.Deployment) error {
//	params, err := json.Marshal([]interface{}{deployment})
//	if err != nil {
//		return err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.CreateDeployment",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return err
//	}
//
//	if rsp.Error != nil {
//		return err
//	}
//
//	return nil
//}
//
//func deleteDeploymentsJsonRPC(url string, deployment ctypes.Deployment) error {
//	params, err := json.Marshal([]interface{}{deployment, false})
//	if err != nil {
//		return err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.CloseDeployment",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return err
//	}
//
//	if rsp.Error != nil {
//		return err
//	}
//
//	return nil
//}
//
//func updateDeploymentsJsonRPC(url string, deployment ctypes.Deployment) error {
//	params, err := json.Marshal([]interface{}{deployment})
//	if err != nil {
//		return err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.UpdateDeployment",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return err
//	}
//
//	if rsp.Error != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getDeploymentLogsJsonRPC(url string, deployment ctypes.Deployment) ([]*ctypes.ServiceLog, error) {
//	params, err := json.Marshal([]interface{}{deployment})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetLogs",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var logs []*ctypes.ServiceLog
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &logs)
//	if err != nil {
//		return nil, err
//	}
//
//	return logs, nil
//}
//
//func getDeploymentEventsJsonRPC(url string, deployment ctypes.Deployment) ([]*ctypes.ServiceEvent, error) {
//	params, err := json.Marshal([]interface{}{deployment})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetEvents",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var event []*ctypes.ServiceEvent
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &event)
//	if err != nil {
//		return nil, err
//	}
//
//	return event, nil
//}
//
//func getDeploymentDomainJsonRPC(url string, id ctypes.DeploymentID) ([]*ctypes.DeploymentDomain, error) {
//	params, err := json.Marshal([]interface{}{id})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetDeploymentDomains",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var domains []*ctypes.DeploymentDomain
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &domains)
//	if err != nil {
//		return nil, err
//	}
//
//	return domains, nil
//}
//
//func addDeploymentDomainJsonRPC(url string, id ctypes.DeploymentID, cert *ctypes.Certificate) error {
//	params, err := json.Marshal([]interface{}{id, cert})
//	if err != nil {
//		return err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.AddDeploymentDomain",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return err
//	}
//
//	var domains []*ctypes.DeploymentDomain
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return err
//	}
//
//	err = json.Unmarshal(b, &domains)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func deleteDeploymentDomainJsonRPC(url string, id ctypes.DeploymentID, host string) error {
//	params, err := json.Marshal([]interface{}{id, host})
//	if err != nil {
//		return err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.DeleteDeploymentDomain",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return err
//	}
//
//	var domains []*ctypes.DeploymentDomain
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return err
//	}
//
//	err = json.Unmarshal(b, &domains)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getDeploymentShellJsonRPC(url string, id ctypes.DeploymentID) (*ctypes.LeaseEndpoint, error) {
//	params, err := json.Marshal([]interface{}{id})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetLeaseShellEndpoint",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var endpoint ctypes.LeaseEndpoint
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &endpoint)
//	if err != nil {
//		return nil, err
//	}
//
//	return &endpoint, nil
//}
//
//func getIngressJsonRPC(url string, id ctypes.DeploymentID) (*ctypes.Ingress, error) {
//	params, err := json.Marshal([]interface{}{id})
//	if err != nil {
//		return nil, err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.GetIngress",
//		Params:  params,
//		ID:      1,
//	}
//
//	rsp, err := requestJsonRPC(url, req)
//	if err != nil {
//		return nil, err
//	}
//
//	var ingress ctypes.Ingress
//	b, err := json.Marshal(rsp.Result)
//	if err != nil {
//		return nil, err
//	}
//
//	err = json.Unmarshal(b, &ingress)
//	if err != nil {
//		return nil, err
//	}
//
//	return &ingress, err
//}
//
//func updateIngressJsonRPC(url string, id ctypes.DeploymentID, annotations map[string]string) error {
//	params, err := json.Marshal([]interface{}{id, annotations})
//	if err != nil {
//		return err
//	}
//
//	req := model.LotusRequest{
//		Jsonrpc: "2.0",
//		Method:  "titan.UpdateIngress",
//		Params:  params,
//		ID:      1,
//	}
//
//	_, err = requestJsonRPC(url, req)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func requestJsonRPC(url string, req model.LotusRequest) (*model.LotusResponse, error) {
//	jsonData, err := json.Marshal(req)
//	if err != nil {
//		return nil, err
//	}
//	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
//	if err != nil {
//		return nil, err
//	}
//
//	//token := config.Cfg.ContainerManager.Token
//	request.Header.Add("Content-Type", "application/json")
//	request.Header.Add("Authorization", "Bearer "+token)
//	resp, err := http.DefaultClient.Do(request)
//	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
//	//if err != nil {
//	//	return nil, err
//	//}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println(string(body))
//
//	var rsp model.LotusResponse
//	err = json.Unmarshal(body, &rsp)
//	if err != nil {
//		return nil, err
//	}
//
//	if rsp.Error != nil {
//		return nil, xerrors.New(rsp.Error.Message)
//	}
//
//	return &rsp, nil
//}
