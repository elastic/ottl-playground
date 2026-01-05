//go:build !ottl_ptr

/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. Elasticsearch B.V. licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package internal

import (
	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlresource"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlscope"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspanevent"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

// Validation functions for OTTL versions before v0.142.0 (non-pointer TransformContext types)

func validateLogStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottllog.TransformContext]()
	parser, err := ottllog.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateSpanStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlspan.TransformContext]()
	parser, err := ottlspan.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateSpanEventStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlspanevent.TransformContext]()
	parser, err := ottlspanevent.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateMetricStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlmetric.TransformContext]()
	parser, err := ottlmetric.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateDataPointStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottldatapoint.TransformContext]()
	parser, err := ottldatapoint.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateResourceStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlresource.TransformContext]()
	parser, err := ottlresource.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateScopeStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlscope.TransformContext]()
	parser, err := ottlscope.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}
