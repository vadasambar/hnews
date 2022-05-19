/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/vadasambar/hnews/api/v1"
)

// HNewsReconciler reconciles a HNews object
type HNewsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	scoreRegex = regexp.MustCompile("^[[:space:]]*(>|>=|=|!=|<|<=)[[:space:]]*([[:digit:]]+[[:space:]]*$)")
)

type Type string

// Based on the possible types specified in https://github.com/HackerNews/API#items
const (
	Story   Type = "story"
	Comment Type = "comment"
	Job     Type = "job"
	Poll    Type = "poll"
	PollOpt Type = "pollopt"
)

// Generated using https://mholt.github.io/json-to-go/
// by converting get Id http json response to go struct
type GetIdResponse struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        Type   `json:"type"`
	URL         string `json:"url"`
}

type Filter struct {
	Limit       int        `json:"limit"`
	Type        string     `json:"type"`
	Score       Comparison `json:"score"`
	Descendants Comparison `json:"descendents"`
}

type Comparison string

//+kubebuilder:rbac:groups=apps.vadasambar.com,resources=hnews,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.vadasambar.com,resources=hnews/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.vadasambar.com,resources=hnews/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HNews object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *HNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var hn appsv1.HNews
	err := r.Client.Get(ctx, req.NamespacedName, &hn)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Log.Info("unable to fetch hnews k8s resource", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, nil
		}
		log.Log.Error(err, "unable to fetch hnews k8s resource", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// RESUME HERE
	// defaultFiltersSet := true
	// if hn.Spec.Filter.Type == "" || hn.Spec.Filter.Limit == 0 || hn.Spec.Filter.Score == "" || hn.Spec.Filter.Descendants == "" {
	// 	defaultFiltersSet = false
	// }

	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		log.Log.Error(err, "error getting response from /topstories.json API")
		return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 30}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Log.Error(err, "error reading response from /topstories.json API")
		return ctrl.Result{}, err
	}
	var ids []json.Number
	err = json.Unmarshal(body, &ids)
	if err != nil {
		log.Log.Error(err, "error unmarshalling /topstories.json API response")
		return ctrl.Result{}, err
	}

	filter := Filter{
		Limit:       10,
		Score:       ">=100",
		Descendants: ">3",
	}

	result := []GetIdResponse{}

	count := 0
	for _, id := range ids {
		if filter.Limit == count {
			break
		}
		resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json", id))
		if err != nil {
			log.Log.Error(err, "error getting response from /item/{item-id}.json API")
			return ctrl.Result{Requeue: true, RequeueAfter: time.Second * 30}, err
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Log.Error(err, "error reading response from /item/{item-id}.json API")
			return ctrl.Result{}, err
		}
		var getIdResp GetIdResponse
		err = json.Unmarshal(body, &getIdResp)
		if err != nil {
			log.Log.Error(err, "error unmarshalling /item/{item-id}.json API response")
			return ctrl.Result{}, err
		}

		if evalCond(getIdResp.Score, filter.Score) && count < filter.Limit && evalCond(getIdResp.Descendants, filter.Descendants) {
			result = append(result, getIdResp)
			count++
		}
	}

	for _, r := range result {
		fmt.Println("TITLE:", r.Title, "SCORE", r.Score, "DESCENDENTS", r.Descendants)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.HNews{}).
		Complete(r)
}

func evalCond(value int, cond Comparison) bool {
	// https: //play.golang.com/p/B8ZgghEBK4k

	result := scoreRegex.FindAllStringSubmatch(string(cond), -1)
	if len(result[0]) < 3 {
		return false
	}

	comparisonOperator := strings.TrimSpace(result[0][1])
	condValue, err := strconv.Atoi(strings.TrimSpace(result[0][2]))
	if err != nil {
		fmt.Println("err", err)
		return false
	}
	switch comparisonOperator {
	case ">":
		return value > condValue
	case ">=":
		return value >= condValue
	case "<":
		return value < condValue
	case "<=":
		return value <= condValue
	case "=":
		return value == condValue
	case "!=":
		return value != condValue
	}

	return false
}
