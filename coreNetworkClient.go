package libDatabox

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	dockerNetworkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

type CoreNetworkClient struct {
	cli     *client.Client
	request *http.Client
	CM_KEY  string
}

type NetworkConfig struct {
	NetworkName string
	DNS         string
}

type PostNetworkConfig struct {
	NetworkName string
	IPv4Address string
}

func NewCoreNetworkClient(containerManagerKeyPath string, request *http.Client) CoreNetworkClient {

	cli, _ := client.NewEnvClient()

	cmKeyBytes, err := ioutil.ReadFile(containerManagerKeyPath)
	var cmKey string
	if err != nil {
		Warn("Warning:: failed to read core-network CM_KEY using empty string")
		cmKey = ""
	} else {
		cmKey = b64.StdEncoding.EncodeToString([]byte(cmKeyBytes))
	}

	return CoreNetworkClient{
		cli:     cli,
		request: request,
		CM_KEY:  cmKey,
	}
}

func (cnc CoreNetworkClient) PreConfig(localContainerName string, sla SLA) NetworkConfig {

	networkName := localContainerName + "-network"

	internal := true
	if sla.DataboxType == "driver" {
		internal = false
	}

	//check for an existing network
	f := filters.NewArgs()
	f.Add("name", networkName)
	networkList, _ := cnc.cli.NetworkList(context.Background(), types.NetworkListOptions{Filters: f})

	var network types.NetworkResource
	var err error

	if len(networkList) > 0 {
		//network exists
		network, err = cnc.cli.NetworkInspect(context.Background(), networkList[0].ID, types.NetworkInspectOptions{})
		if err != nil {
			Err("[PreConfig] NetworkInspect1 Error " + err.Error())
		}
		Debug("[PreConfig] using existing network " + network.Name)
	} else {
		//create network
		networkCreateResponse, err := cnc.cli.NetworkCreate(context.Background(), networkName, types.NetworkCreate{
			Internal:   internal,
			Driver:     "overlay",
			Attachable: true,
		})
		if err != nil {
			Err("[PreConfig] NetworkCreate Error " + err.Error())
		}

		network, err = cnc.cli.NetworkInspect(context.Background(), networkCreateResponse.ID, types.NetworkInspectOptions{})
		if err != nil {
			Err("[PreConfig] NetworkInspect2 Error " + err.Error())
		}

		//find core network
		f := filters.NewArgs()
		f.Add("name", "databox-network") //TODO hardcoded
		coreNetwork, err := cnc.cli.ContainerList(context.Background(), types.ContainerListOptions{
			Filters: f,
		})
		if err != nil {
			Err("[PreConfig] ContainerList Error " + err.Error())
		}

		//attach to core-network
		err = cnc.cli.NetworkConnect(
			context.Background(),
			network.ID,
			coreNetwork[0].ID,
			&dockerNetworkTypes.EndpointSettings{},
		)
		if err != nil {
			Err("[PreConfig] NetworkConnect Error " + err.Error())
		}
		time.Sleep(time.Second * 5)
		//refresh network status
		network, err = cnc.cli.NetworkInspect(context.Background(), networkCreateResponse.ID, types.NetworkInspectOptions{})
		if err != nil {
			Err("[PreConfig] NetworkInspect3 Error " + err.Error())
		}
	}

	//find core-network IP on the new network to used as DNS
	var ipOnNewNet string
	for _, cont := range network.Containers {
		Debug("contName=" + cont.Name)
		if cont.Name == "databox-network" {
			ipOnNewNet = strings.Split(cont.IPv4Address, "/")[0]
			break
		}
	}

	Debug("[PreConfig]" + networkName + " " + ipOnNewNet)

	return NetworkConfig{NetworkName: networkName, DNS: ipOnNewNet}
}

func (cnc CoreNetworkClient) NetworkOfService(service swarm.Service, serviceName string) (PostNetworkConfig, error) {
	fmt.Println("NetworkOfService")

	netConfig := PostNetworkConfig{}

	netConfig.NetworkName = serviceName + "-network"
	//get IP of service
	netFilters := filters.NewArgs()
	netFilters.Add("name", netConfig.NetworkName)
	networks, err := cnc.cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: netFilters,
	})
	ChkErr(err)

	if len(networks) < 1 {
		fmt.Println("Can't find network")

		return netConfig, errors.New("Can't find network " + netConfig.NetworkName)
	}

	for _, net := range networks {
		fmt.Println("network name", net.Name)
		netInfo, _ := cnc.cli.NetworkInspect(context.Background(), net.ID, types.NetworkInspectOptions{})
		for _, endpoint := range netInfo.Containers {
			fmt.Println(endpoint.Name, endpoint.IPv4Address)
			if cnc.toServiceName(endpoint.Name) == serviceName {
				//				netConfig.IPv4Address = strings.Split(endpoint.IPv4Address, "/")[0]
				netConfig.IPv4Address = endpoint.IPv4Address
				break
			}
		}
	}

	fmt.Println("returning ", netConfig)
	return netConfig, nil

}

func (cnc CoreNetworkClient) toServiceName(containerName string) string {

	parts := strings.Split(containerName, ".")

	return parts[0]
}

func (cnc CoreNetworkClient) PostUninstall(name string, netConfig PostNetworkConfig) error {

	return cnc.DisconnectEndpoints(name, netConfig)

	//TODO remove empty networks !!!
}

func (cnc CoreNetworkClient) post(LogFnName string, data []byte, URL string) error {
	Debug("[CoreNetworkClient." + LogFnName + "] POSTED JSON :: " + string(data))
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
	if err != nil {
		Err("[" + LogFnName + "] Error:: " + err.Error())
		return err
	}
	req.Header.Set("x-api-key", cnc.CM_KEY)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	resp, err := cnc.request.Do(req)

	if err != nil {
		Err("[" + LogFnName + "] Error " + err.Error())
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		response, _ := ioutil.ReadAll(resp.Body)
		Err("[" + LogFnName + "] PostError StatusCode=" + strconv.Itoa(resp.StatusCode) + " data=" + string(data) + "response=" + string(response))
		return err
	}

	return nil
}

func (cnc CoreNetworkClient) ConnectEndpoints(serviceName string, peers []string) error {

	type postData struct {
		Name  string   `json:"name"`
		Peers []string `json:"peers"`
	}

	data := postData{
		Name:  serviceName,
		Peers: peers,
	}

	postBytes, _ := json.Marshal(data)

	return cnc.post("ConnectEndpoints", postBytes, "https://databox-network:8080/connect")
}

func (cnc CoreNetworkClient) DisconnectEndpoints(serviceName string, netConfig PostNetworkConfig) error {

	type postData struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	}

	data := postData{
		Name: serviceName,
		IP:   netConfig.IPv4Address,
	}

	postBytes, _ := json.Marshal(data)

	return cnc.post("DisconnectEndpoints", postBytes, "https://databox-network:8080/disconnect")
}

func (cnc CoreNetworkClient) RegisterPrivileged() error {

	cmIP, err := cnc.getCmIP()
	if err != nil {
		return err
	}

	jsonStr := "{\"src_ip\":\"" + cmIP + "\"}"
	return cnc.post("RegisterPrivileged", []byte(jsonStr), "https://databox-network:8080/privileged")

}

func (cnc CoreNetworkClient) ServiceRestart(serviceName string, oldIP string, newIP string) error {

	type postData struct {
		Name  string `json:"name"`
		OldIP string `json:"old_ip"`
		NewIP string `json:"new_ip"`
	}

	data := postData{
		Name:  serviceName,
		OldIP: oldIP,
		NewIP: newIP,
	}
	postBytes, _ := json.Marshal(data)
	return cnc.post("ServiceRestart", postBytes, "https://databox-network:8080/restart")

}

func (cnc CoreNetworkClient) getCmIP() (string, error) {

	f := filters.NewArgs()
	f.Add("name", "container-manager")

	containerList, _ := cnc.cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: f,
	})

	if len(containerList) < 1 {
		Err("[getCmIP] Error no CM found for core-network")
		return "", errors.New("No CM found for core-network")
	}

	if _, ok := containerList[0].NetworkSettings.Networks["databox-system-net"]; ok {
		return containerList[0].NetworkSettings.Networks["databox-system-net"].IPAddress, nil
	}

	Err("[getCmIP] CM not on core-network")
	return "", errors.New("CM not on core-network")

}
