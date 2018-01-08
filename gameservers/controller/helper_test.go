// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"time"

	"context"

	agonfake "github.com/agonio/agon/pkg/client/clientset/versioned/fake"
	"github.com/agonio/agon/pkg/client/informers/externalversions"
	"github.com/sirupsen/logrus"
	extfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/client-go/informers"
	kubefake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

// holder for all my fakes and mocks
type mocks struct {
	kubeClient             *kubefake.Clientset
	kubeInformationFactory informers.SharedInformerFactory
	extClient              *extfake.Clientset
	agonClient             *agonfake.Clientset
	agonInformerFactory    externalversions.SharedInformerFactory
}

func newMocks() mocks {
	kubeClient := &kubefake.Clientset{}
	kubeInformationFactory := informers.NewSharedInformerFactory(kubeClient, 30*time.Second)
	extClient := &extfake.Clientset{}
	agonClient := &agonfake.Clientset{}
	agonInformerFactory := externalversions.NewSharedInformerFactory(agonClient, 30*time.Second)
	m := mocks{
		kubeClient:             kubeClient,
		kubeInformationFactory: kubeInformationFactory,
		extClient:              extClient,
		agonClient:             agonClient,
		agonInformerFactory:    agonInformerFactory}
	return m
}

func startInformers(mocks mocks, sync ...cache.InformerSynced) (<-chan struct{}, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	stop := ctx.Done()

	mocks.kubeInformationFactory.Start(stop)
	mocks.agonInformerFactory.Start(stop)

	logrus.Info("Wait for cache sync")
	if !cache.WaitForCacheSync(stop, sync...) {
		panic("Cache never synced")
	}

	return stop, cancel
}