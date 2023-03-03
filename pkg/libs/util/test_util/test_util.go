package test_util

/*
Certificate tests

Files _test_cert.pem_ and _test_key.pem_ have been created with following command:

> openssl req -x509 -nodes -batch -newkey rsa:4096 -keyout test_key.pem -out test_cert.pem -sha256 -days 365
*/

import (
	_ "embed"
)

//go:embed test_cert.pem
var TestCertPEM string

//go:embed test_key.pem
var TestKeyPEM string

const TestCertExpiryDate string = "2024-03-02T13:29:09Z"
