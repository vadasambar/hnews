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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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

// HNewsSpec defines the desired state of HNews
type HNewsSpec struct {
	Filter Filter `json:"filter"`
}

// HNewsStatus defines the observed state of HNews
type HNewsStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Links []Link `json:"link"`
}

type Link struct {
	HNewsUrl    string `json:"hnews_url"`
	ArticleUrl  string `json:"article_url"`
	Descendents int    `json:"descendents"`
	Score       int    `json:"score"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HNews is the Schema for the hnews API
type HNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HNewsSpec   `json:"spec,omitempty"`
	Status HNewsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HNewsList contains a list of HNews
type HNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HNews{}, &HNewsList{})
}
