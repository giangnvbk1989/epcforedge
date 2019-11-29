// Copyright 2019 Intel Corporation, Inc. All rights reserved
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

package ngcaf

import (
	"context"
	"net/http"
	"strconv"
)

func deleteSubscription(cliCtx context.Context, afCtx *afContext,
	sID string) (*http.Response, error) {

	cliCfg := NewConfiguration(afCtx)
	cli := NewClient(cliCfg)

	resp, err := cli.TrafficInfluSubDeleteAPI.SubscriptionDelete(cliCtx,
		afCtx.cfg.AfID, sID)

	if err != nil {
		return nil, err
	}
	return resp, nil

}

// DeleteSubscription function
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		resp           *http.Response
		subscriptionID string
	)

	afCtx := r.Context().Value(keyType("af-ctx")).(*afContext)
	cliCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	subscriptionID, err = getSubsIDFromURL(r.URL)
	if err != nil {
		log.Errf("Traffic Influence Subscription delete: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp, err = deleteSubscription(cliCtx, afCtx, subscriptionID)
	if err != nil {
		log.Errf("Traffic Influence Subscription delete: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if interMap, ok := afCtx.subscriptions[subscriptionID]; ok {
		for transID := range interMap {
			var i int
			if i, err = strconv.Atoi(transID); err != nil {
				log.Errf("Error converting transID to integer: %v", err)
			} else {
				delete(afCtx.transactions, i)
				log.Infof("Deleted transaction ID %v", i)
			}
		}
		delete(afCtx.subscriptions, subscriptionID)
	}

	w.WriteHeader(resp.StatusCode)
}
