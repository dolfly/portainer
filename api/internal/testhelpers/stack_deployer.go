package testhelpers

import portainer "github.com/portainer/portainer/api"

type TestStackDeployer struct {
	DeployComposeCallCount int
	DeploySwarmCallCount   int
	LastPrune              bool
}

func NewTestStackDeployer() *TestStackDeployer {
	return &TestStackDeployer{}
}

func (d *TestStackDeployer) DeployComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, prune, forcePullImage, forceRecreate bool) error {
	d.DeployComposeCallCount++
	d.LastPrune = prune
	return nil
}

func (d *TestStackDeployer) DeploySwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, prune, pullImage bool) error {
	d.DeploySwarmCallCount++
	d.LastPrune = prune
	return nil
}

func (d *TestStackDeployer) DeployKubernetesStack(stack *portainer.Stack, endpoint *portainer.Endpoint, user *portainer.User) error {
	return nil
}

func (d *TestStackDeployer) DeployRemoteComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, prune, forcePullImage, forceRecreate bool) error {
	return nil
}

func (d *TestStackDeployer) UndeployRemoteComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint) error {
	return nil
}

func (d *TestStackDeployer) StartRemoteComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry) error {
	return nil
}

func (d *TestStackDeployer) StopRemoteComposeStack(stack *portainer.Stack, endpoint *portainer.Endpoint) error {
	return nil
}

func (d *TestStackDeployer) DeployRemoteSwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry, prune, pullImage bool) error {
	return nil
}

func (d *TestStackDeployer) UndeployRemoteSwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint) error {
	return nil
}

func (d *TestStackDeployer) StartRemoteSwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint, registries []portainer.Registry) error {
	return nil
}

func (d *TestStackDeployer) StopRemoteSwarmStack(stack *portainer.Stack, endpoint *portainer.Endpoint) error {
	return nil
}
