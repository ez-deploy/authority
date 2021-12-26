package main

import (
	"io/ioutil"
	"log"
	"net"

	"github.com/ez-deploy/authority/db"
	"github.com/ez-deploy/authority/service"

	pb "github.com/ez-deploy/protobuf/authority"
	"github.com/wuhuizuo/sqlm"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	_ "github.com/go-sql-driver/mysql"
)

const cfgFilename = "authorityCfg.yaml"

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:80")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc, err := newAuthorityServiceFromCfg()
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthorityOpsServer(s, svc)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newAuthorityServiceFromCfg() (*service.Service, error) {
	rawCfg, err := ioutil.ReadFile(cfgFilename)
	if err != nil {
		return nil, err
	}

	configMap := map[string]string{}
	if err := yaml.Unmarshal(rawCfg, configMap); err != nil {
		return nil, err
	}

	database := sqlm.Database{
		Driver: sqlm.DriverMysql,
		DSN:    configMap["dsn"],
	}
	if err := database.Create(); err != nil {
		return nil, err
	}

	authorityTable := &sqlm.Table{
		Database:  &database,
		TableName: "authority",
	}
	authorityTable.SetRowModel(db.AuthorityRawModel)
	if err := authorityTable.Create(); err != nil {
		return nil, err
	}

	return &service.Service{AuthorityTable: authorityTable}, nil
}
