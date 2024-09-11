package main

import (
	"context"
	"fmt"
	"github.com/gnasnik/titan-container-api/api"
	"github.com/gnasnik/titan-container-api/config"
	"github.com/gnasnik/titan-container-api/core/dao"
	"github.com/gnasnik/titan-container-api/core/generated/model"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("reading config file: %v\n", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshaling config file: %v\n", err)
	}

	if err := dao.Init(&cfg); err != nil {
		log.Fatalf("initital: %v\n", err)
	}

	ec, err := api.NewEtcdClient(cfg.EtcdUser, cfg.EtcdPassword, cfg.EtcdAddresses)
	if err != nil {
		log.Fatal(err)
	}

	schedulers, err := api.FetchSchedulersFromEtcd(context.Background(), ec)
	if err != nil {
		log.Fatal(err)
	}

	id := "c_0131ee95-5a43-43e4-8e79-b73b12f78b0e"

	_, providers, err := dao.GetProvidersWithResource(context.Background(), "", model.QueryOption{ID: id})
	if err != nil {
		log.Fatal(err)
	}

	for _, provider := range providers {
		//provider.AreaID

		for _, scheduler := range schedulers {
			if scheduler.AreaId != provider.AreaID {
				continue
			}

			stat, err := scheduler.Api.GetStatistics(context.Background(), id)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%+v", stat)
		}

	}

	log.Println("Success")
}
