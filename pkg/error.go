package istioapi

import "errors"

var (
	ErrorListIstioResource   = errors.New("error to list istio resources")
	ErrorGetIstioResource    = errors.New("error to get istio resource")
	ErrorCreateIstioResource = errors.New("error to create istio resource")
	ErrorUpdateIstioResource = errors.New("error to update istio resource")
	ErrorDeleteIstioResource = errors.New("error to delete istio resource")
)
