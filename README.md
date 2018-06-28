# eosmonitor

EOS BP monitor with TICK Stack.

## Features
- [x] Use influxdata's TICK Stack.
- [x] Gather [get_info](https://developers.eos.io/eosio-nodeos/reference#get_info) from every nodeos's localhost RPC interface.
- [x] Use TICKScript to alert when the latest `head_block_time` is more than 10 seconds earlier `now()`. It usually means that the BP stopped producing blocks.
