// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// openpitrix api gateway
package main

import (
	"openpitrix.io/openpitrix/pkg/apigateway"
	"openpitrix.io/openpitrix/pkg/config"
)

func main() {

	cfg := config.LoadConf()

	apigateway.Serve(cfg)
}
