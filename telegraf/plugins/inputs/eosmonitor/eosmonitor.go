// Package eosmonitor
package eosmonitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

const sampleConfig = `
  ## Set the EOS rpc_url
  rpc_url = "http://localhost:8888/v1/chain/get_info"
`

type EosMonitor struct {
	RpcUrl string
}

func (em *EosMonitor) SampleConfig() string {
	return sampleConfig
}

func (em *EosMonitor) Description() string {
	return "EOS Monitor check chain status from RPC interface"
}

func (em *EosMonitor) Gather(acc telegraf.Accumulator) error {

	resp, err := http.Get(em.RpcUrl)
	if err != nil {
		acc.AddError(fmt.Errorf("EOS monitor http get request failed, err: %v", err))
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		acc.AddError(fmt.Errorf("EOS monitor http get body failed, err: %v", err))
		return nil
	}

	//serverVersion := gjson.GetBytes(body, "server_version").String()
	//chainId := gjson.GetBytes(body, "chain_id").String()
	headBlockNum := gjson.GetBytes(body, "head_block_num").Int()
	lastIrreversibleBlockNum := gjson.GetBytes(body, "last_irreversible_block_num").Int()
	//lastIrreversibleBlockId := gjson.GetBytes(body, "last_irreversible_block_id").String()
	//headBlockId := gjson.GetBytes(body, "head_block_id").String()
	headBlockTime := cast.ToTime(gjson.GetBytes(body, "head_block_time").String()).Unix()
	//headBlockProducer := gjson.GetBytes(body, "head_block_producer").String()

	fields := make(map[string]interface{})

	//fields["server_version"] = serverVersion
	//fields["chain_id"] = chainId
	fields["head_block_num"] = headBlockNum
	fields["last_irreversible_block_num"] = lastIrreversibleBlockNum
	//fields["last_irreversible_block_id"] = lastIrreversibleBlockId
	//fields["head_block_id"] = headBlockId
	fields["head_block_time"] = headBlockTime
	//fields["head_block_producer"] = headBlockProducer
	fields["now"] = time.Now().Unix()

	//fmt.Printf("eosmonitor, fields: %+v\n", fields)

	tags := make(map[string]string)

	acc.AddFields("eosmonitor", fields, tags)

	return nil
}

func NewEosMonitor() *EosMonitor {
	return &EosMonitor{
		RpcUrl: "http://localhost:8888/v1/chain/get_info",
	}
}

func init() {
	inputs.Add("eosmonitor", func() telegraf.Input {
		return NewEosMonitor()
	})
}
