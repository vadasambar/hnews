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
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "github.com/vadasambar/hnews/api/v1"
	helpers "github.com/vadasambar/hnews/pkg/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HNewsReconciler reconciles a HNews object
type HNewsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	defaultDescendents = ">5"
	defaultScore       = ">200"
	defaultLimit       = 5
	defaultType        = string(appsv1.Story)
	topStoriesUrl      = "https://hacker-news.firebaseio.com/v0/topstories.json"
	getIdRespUrl       = "https://hacker-news.firebaseio.com/v0/item/%s.json"
	hnewsArticleUrl    = "https://news.ycombinator.com/item?id=%d"
)

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

	if hn.Spec.Filter.Type == "" || hn.Spec.Filter.Limit == 0 || hn.Spec.Filter.Score == "" || hn.Spec.Filter.Descendants == "" {
		if hn.Spec.Filter.Type == "" {
			hn.Spec.Filter.Type = defaultType
		}

		if hn.Spec.Filter.Limit == 0 {
			hn.Spec.Filter.Limit = defaultLimit
		}

		if hn.Spec.Filter.Score == "" {
			hn.Spec.Filter.Score = defaultScore
		}

		if hn.Spec.Filter.Descendants == "" {
			hn.Spec.Filter.Descendants = defaultDescendents
		}

		if err := r.Update(ctx, &hn); err != nil {
			log.Log.Error(err, "unable to update hnews", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{RequeueAfter: time.Second * 30}, err
		}
		// reconcile is triggered automatically if the spec is updated
		// no need to use `Requeue` below`
		return ctrl.Result{}, nil
	}

	resp, err := http.Get(topStoriesUrl)
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

	hn.Status.Links = []appsv1.Link{}
	count := 0
	for _, id := range ids {
		if hn.Spec.Filter.Limit == count {
			break
		}
		resp, err := http.Get(fmt.Sprintf(getIdRespUrl, id))
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
		var getIdResp appsv1.GetIdResponse
		err = json.Unmarshal(body, &getIdResp)
		if err != nil {
			log.Log.Error(err, "error unmarshalling /item/{item-id}.json API response")
			return ctrl.Result{}, err
		}

		if helpers.EvalCond(getIdResp.Score, hn.Spec.Filter.Score) && count < hn.Spec.Filter.Limit && helpers.EvalCond(getIdResp.Descendants, hn.Spec.Filter.Descendants) {
			hn.Status.Links = append(hn.Status.Links, appsv1.Link{
				HNewsUrl:    fmt.Sprintf(hnewsArticleUrl, getIdResp.ID),
				ArticleUrl:  getIdResp.URL,
				Descendents: getIdResp.Descendants,
				Score:       getIdResp.Score,
			})
			count++
		}
	}

	hn.Status.LastSyncedAt = metav1.NewTime(time.Now())
	if err := r.Status().Update(ctx, &hn); err != nil {
		log.Log.Error(err, "unable to update hnews status", "name", req.Name, "namespace", req.Namespace)
		return ctrl.Result{RequeueAfter: time.Second * 30}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.HNews{}).
		Complete(r)
}
