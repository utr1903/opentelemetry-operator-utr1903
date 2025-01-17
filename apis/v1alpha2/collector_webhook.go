// Copyright The OpenTelemetry Authors
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

package v1alpha2

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	_ admission.CustomValidator = &CollectorWebhook{}
	_ admission.CustomDefaulter = &CollectorWebhook{}
)

// +kubebuilder:object:generate=false
type CollectorWebhook struct {
}

func (c CollectorWebhook) Default(_ context.Context, obj runtime.Object) error {
	otelcol, ok := obj.(*OpenTelemetryCollector)
	if !ok {
		return fmt.Errorf("expected an OpenTelemetryCollector, received %T", obj)
	}
	return c.defaulter(otelcol)
}

func (c CollectorWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	otelcol, ok := obj.(*OpenTelemetryCollector)
	if !ok {
		return nil, fmt.Errorf("expected an OpenTelemetryCollector, received %T", obj)
	}
	return c.validate(otelcol)
}

func (c CollectorWebhook) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	otelcol, ok := newObj.(*OpenTelemetryCollector)
	if !ok {
		return nil, fmt.Errorf("expected an OpenTelemetryCollector, received %T", newObj)
	}
	return c.validate(otelcol)
}

func (c CollectorWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	otelcol, ok := obj.(*OpenTelemetryCollector)
	if !ok || otelcol == nil {
		return nil, fmt.Errorf("expected an OpenTelemetryCollector, received %T", obj)
	}
	return c.validate(otelcol)
}

func (c CollectorWebhook) defaulter(r *OpenTelemetryCollector) error {
	return nil
}

func (c CollectorWebhook) validate(r *OpenTelemetryCollector) (admission.Warnings, error) {
	warnings := admission.Warnings{}

	nullObjects := r.Spec.Config.nullObjects()
	if len(nullObjects) > 0 {
		warnings = append(warnings, fmt.Sprintf("Collector config spec.config has null objects: %s. For compatibility tooling (kustomize and kubectl edit) it is recommended to use empty obejects e.g. batch: {}.", strings.Join(nullObjects, ", ")))
	}
	return warnings, nil
}
