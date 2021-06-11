package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/juxuny/env"
	pb "github.com/juxuny/supervisor/proxy"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedProxyServer
}

func (s *server) Status(ctx context.Context, in *pb.StatusReq) (*pb.StatusResp, error) {
	var resp pb.StatusResp
	currentStatus, _ := proxyServer.Status()
	resp.Status = &pb.Status{
		ListenPort: currentStatus.ListenPort,
		Remote:     currentStatus.Remote,
	}
	return &resp, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateReq) (*pb.UpdateResp, error) {
	var resp pb.UpdateResp
	if err := proxyServer.UpdateRemote(in.Status.Remote); err != nil {
		return nil, errors.Wrap(err, "update failed")
	}
	currentStatus, _ := proxyServer.Status()
	resp.Status = &pb.Status{
		ListenPort: currentStatus.ListenPort,
		Remote:     currentStatus.Remote,
	}
	return &resp, nil
}

func (s *server) Check(ctx context.Context, req *pb.CheckReq) (*pb.CheckResp, error) {
	var resp pb.CheckResp
	var checkErr error
	if req.Type == pb.HealthCheckType_TypeDefault {
		u := fmt.Sprintf("http://%s:%d%s", req.Host, req.Port, req.Path)
		fmt.Println("checking:", u)
		checkErr = checkHttp(u)
	} else if req.Type == pb.HealthCheckType_TypeTcp {
		u := fmt.Sprintf("%s:%d", req.Host, req.Port)
		fmt.Println("checking:", u)
		checkErr = checkTcp(u)
	} else {
		return nil, errors.Errorf("unknown HealthCheckType:%v", req.Type)
	}
	if checkErr != nil {
		fmt.Println(checkErr)
		return nil, checkErr
	}
	resp.Code = 0
	resp.Msg = "success"
	return &resp, nil
}

var (
	configFile  string
	fromEnv     bool
	proxyServer pb.IServer
)

func main() {
	flag.StringVar(&configFile, "c", "config/proxy.yaml", "proxy config yaml")
	flag.BoolVar(&fromEnv, "e", false, "use config from environment variable")
	flag.Parse()
	var err error
	var proxyConfig *pb.Config
	if fromEnv {
		proxyConfig = &pb.Config{Proxy: pb.Proxy{
			ControlPort: uint32(env.GetInt("CONTROL_PORT", 50050)),
			ListenPort:  uint32(env.GetInt("LISTEN_PORT", 8888)),
			Remote:      env.GetString("REMOTE", "127.0.0.1:8080"),
		}}
	} else {
		proxyConfig, err = pb.Parse(configFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	proxyServer = pb.NewServer(proxyConfig.Proxy)
	go proxyServer.Start()
	addr := fmt.Sprintf(":%d", int(proxyConfig.Proxy.ControlPort))
	fmt.Println("control plane listen ", addr)
	fmt.Println("proxy port: ", proxyConfig.Proxy.ListenPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProxyServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
