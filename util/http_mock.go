// Automatically generated by MockGen. DO NOT EDIT!
// Source: http.go

package util

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	time "time"
)

// Mock of HttpClientInterface interface
type MockHttpClientInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockHttpClientInterfaceRecorder
}

// Recorder for MockHttpClientInterface (not exported)
type _MockHttpClientInterfaceRecorder struct {
	mock *MockHttpClientInterface
}

func NewMockHttpClientInterface(ctrl *gomock.Controller) *MockHttpClientInterface {
	mock := &MockHttpClientInterface{ctrl: ctrl}
	mock.recorder = &_MockHttpClientInterfaceRecorder{mock}
	return mock
}

func (_m *MockHttpClientInterface) EXPECT() *_MockHttpClientInterfaceRecorder {
	return _m.recorder
}

func (_m *MockHttpClientInterface) DoGet(url string, unmarshalledResponse interface{}) (int, error) {
	ret := _m.ctrl.Call(_m, "DoGet", url, unmarshalledResponse)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoGet(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoGet", arg0, arg1)
}

func (_m *MockHttpClientInterface) DoPost(url string, body interface{}, unmarshalledResponse interface{}) (int, error) {
	ret := _m.ctrl.Call(_m, "DoPost", url, body, unmarshalledResponse)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoPost(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoPost", arg0, arg1, arg2)
}

func (_m *MockHttpClientInterface) DoPut(url string, body interface{}, unmarshalledResponse interface{}) (int, error) {
	ret := _m.ctrl.Call(_m, "DoPut", url, body, unmarshalledResponse)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoPut(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoPut", arg0, arg1, arg2)
}

func (_m *MockHttpClientInterface) DoDelete(url string, body interface{}, unmarshalledResponse interface{}) (int, error) {
	ret := _m.ctrl.Call(_m, "DoDelete", url, body, unmarshalledResponse)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoDelete(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoDelete", arg0, arg1, arg2)
}

func (_m *MockHttpClientInterface) DoGetWithTimeout(url string, unmarshalledResponse interface{}, timeout time.Duration) (int, error) {
	ret := _m.ctrl.Call(_m, "DoGetWithTimeout", url, unmarshalledResponse, timeout)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoGetWithTimeout(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoGetWithTimeout", arg0, arg1, arg2)
}

func (_m *MockHttpClientInterface) DoPostWithTimeout(url string, body interface{}, unmarshalledResponse interface{}, timout time.Duration) (int, error) {
	ret := _m.ctrl.Call(_m, "DoPostWithTimeout", url, body, unmarshalledResponse, timout)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoPostWithTimeout(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoPostWithTimeout", arg0, arg1, arg2, arg3)
}

func (_m *MockHttpClientInterface) DoRequest(request *http.Request, unmarshalledResponse interface{}, timeout time.Duration) (int, error) {
	ret := _m.ctrl.Call(_m, "DoRequest", request, unmarshalledResponse, timeout)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockHttpClientInterfaceRecorder) DoRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DoRequest", arg0, arg1, arg2)
}

func (_m *MockHttpClientInterface) AddHeader(key string, value string) {
	_m.ctrl.Call(_m, "AddHeader", key, value)
}

func (_mr *_MockHttpClientInterfaceRecorder) AddHeader(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddHeader", arg0, arg1)
}
