package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"io"
	"microservice/models"
	"microservice/transports"
	"net/url"
	"os"
	"time"
)

//客户端直连
//func main2() {
//	tgt, _ := url.Parse("http://127.0.0.1:5678")
//
//	getUserNameEndpoint := httptransport.NewClient("GET", tgt, transports.EncodeUserRequest, transports.DecodeGetUserNameResponse).Endpoint()
//	response, err := getUserNameEndpoint(context.Background(), models.UserRequest{Uid: 100})
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(fmt.Sprintf("%+v", response.(models.UserResponse)))
//
//	delUserEndpoint := httptransport.NewClient("DELETE", tgt, transports.EncodeUserRequest, transports.DecodeDelUserResponse).Endpoint()
//	response, err = delUserEndpoint(context.Background(), models.UserRequest{Uid: 100})
//	if err != nil {
//		panic(err)
//	}
//	if _, ok := response.(models.UserResponse); ok {
//		fmt.Println(fmt.Sprintf("%+v", response.(models.UserResponse)))
//	} else {
//		fmt.Println(fmt.Sprintf("%+v", response.(models.Error)))
//	}
//}

//客户端通过查询consul注册中心获取服务地址（服务发现）
func main() {
	apiclient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		panic(err)
	}

	var (
		sdclient  = consul.NewClient(apiclient)
		logger    = log.NewLogfmtLogger(os.Stdout)
		instancer = consul.NewInstancer(sdclient, logger, "serviceA", []string{"primary"}, true)
	)

	//查询，带负载均衡
	{
		getUserNameFactory := func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
			tgt, _ := url.Parse("http://" + serviceUrl)
			fmt.Println("----------", serviceUrl)
			return httptransport.NewClient("GET", tgt, transports.EncodeUserRequest, transports.DecodeGetUserNameResponse).Endpoint(), nil, nil
		}
		getUserNameEndpointer := sd.NewEndpointer(instancer, getUserNameFactory, logger)

		balancer := lb.NewRoundRobin(getUserNameEndpointer)

		for {
			getUserNameEndpoint, _ := balancer.Endpoint()
			response, err := getUserNameEndpoint(context.Background(), models.UserRequest{Uid: 100})
			if err != nil {
				panic(err)
			}
			fmt.Println(fmt.Sprintf("%+v", response.(models.UserResponse)))

			time.Sleep(time.Second)
		}
	}

	//查询，不带负载均衡
	//{
	//	getUserNameFactory := func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
	//		tgt, _ := url.Parse("http://" + serviceUrl)
	//		fmt.Println("----------", serviceUrl)
	//		return httptransport.NewClient("GET", tgt, transports.EncodeUserRequest, transports.DecodeGetUserNameResponse).Endpoint(), nil, nil
	//	}
	//	getUserNameEndpointer := sd.NewEndpointer(instancer, getUserNameFactory, logger)
	//
	//	getUserNameEndpoints, _ := getUserNameEndpointer.Endpoints()
	//	getUserNameEndpoint := getUserNameEndpoints[0]
	//
	//	for {
	//		response, err := getUserNameEndpoint(context.Background(), models.UserRequest{Uid: 100})
	//		if err != nil {
	//			panic(err)
	//		}
	//		fmt.Println(fmt.Sprintf("%+v", response.(models.UserResponse)))
	//
	//		time.Sleep(time.Second)
	//	}
	//}

	//删除，不带负载均衡
	//{
	//	delUserFactory := func(serviceUrl string) (endpoint.Endpoint, io.Closer, error) {
	//		tgt, _ := url.Parse("http://" + serviceUrl)
	//		fmt.Println("----------", serviceUrl)
	//		return httptransport.NewClient("DELETE", tgt, transports.EncodeUserRequest, transports.DecodeDelUserResponse).Endpoint(), nil, nil
	//	}
	//	delUserEndpointer := sd.NewEndpointer(instancer, delUserFactory, logger)
	//
	//	delUserEndpoints, _ := delUserEndpointer.Endpoints()
	//	delUserEndpoint := delUserEndpoints[0]
	//
	//	response, err := delUserEndpoint(context.Background(), models.UserRequest{Uid: 100})
	//	if err != nil {
	//		panic(err)
	//	}
	//	if _, ok := response.(models.UserResponse); ok {
	//		fmt.Println(fmt.Sprintf("%+v", response.(models.UserResponse)))
	//	} else {
	//		fmt.Println(fmt.Sprintf("%+v", response.(models.Error)))
	//	}
	//}
}
