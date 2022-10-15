package net

import (
	"fmt"
	"gateway.bojiu.com/pkg/viper"
	"time"
)

var (
	PoolsCollect = make(map[string]GrpcPools, 0)
)

func InitGrpcPools() {
	for _, proxy := range viper.PsCfg {
		go func(p viper.ProxyCfg) {
			for _, url := range p.Addr {
				b := []byte(url)
				u := string(b[7:])

				PoolsCollect[u] = NewGrpcPools(u, p.Maxgrpc)
			}
		}(proxy)
	}

	////每隔1秒激活连接池
	timer := time.NewTicker(1 * time.Second)
	defer timer.Stop()
	go func() {
		for {
			select {
			case <-timer.C:
				for _, p := range PoolsCollect {
					fmt.Println("ActiveDeadLink.........................................................................")
					p.ActiveDeadLink()
				}
			}
		}
	}()
}
