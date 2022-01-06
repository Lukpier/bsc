package tracers

import (
	"context"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

// setTraceCallConfigDefaultTracer sets the default tracer to "callTracerParity" if none set
func setTraceCallConfigDefaultTracer(config *TraceCallConfig) *TraceCallConfig {
	if config == nil {
		config = &TraceCallConfig{}
	}

	if config.Tracer == nil {
		tracer := "callTracerParity"
		config.Tracer = &tracer
	}

	return config
}

// TraceAPI is the collection of Ethereum full node APIs exposed over
// the private debugging endpoint.
type TraceAPI struct {
	debugAPI *API
}

// NewTraceAPI creates a new API definition for the full node-related
// private debug methods of the Ethereum service.
func NewTraceAPI(debugAPI *API) *TraceAPI {
	return &TraceAPI{debugAPI: debugAPI}
}

// decorateResponse applies formatting to trace results if needed.
func decorateResponse(res interface{}, config *TraceConfig) (interface{}, error) {
	if config != nil && config.NestedTraceOutput && config.Tracer != nil {
		return decorateNestedTraceResponse(res, *config.Tracer), nil
	}
	return res, nil
}

// decorateNestedTraceResponse formats trace results the way Parity does.
// Docs: https://openethereum.github.io/JSONRPC-trace-module
// Example:
/*
{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "output": "0x",
    "stateDiff": { ... },
    "trace": [ { ... }, ],
    "vmTrace": { ... }
  }
}
*/
func decorateNestedTraceResponse(res interface{}, tracer string) interface{} {
	out := map[string]interface{}{}
	if tracer == "callTracerParity" {
		out["trace"] = res
	} else if tracer == "stateDiffTracer" {
		out["stateDiff"] = res
	} else {
		return res
	}
	return out
}

// CallMany lets you trace a given eth_call. It collects the structured logs created during the execution of EVM
// if the given transaction was added on top of the provided block and returns them as a JSON object.
// You can provide -2 as a block number to trace on top of the pending block.
func (api *TraceAPI) CallMany(ctx context.Context, txs []ethapi.CallArgs, blockNrOrHash rpc.BlockNumberOrHash, config *TraceCallConfig) (interface{}, error) {
	config = setTraceCallConfigDefaultTracer(config)
	return api.debugAPI.TraceCallMany(ctx, txs, blockNrOrHash, config)
}
