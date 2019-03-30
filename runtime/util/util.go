package util

import (
    "context"
)

func DefaultdecodeRequest(_ context.Context, req interface{}) (interface{}, error) {
    return req, nil
}

func DefaultencodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
    return rsp, nil
}
