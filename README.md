# go-orocrm-wsse [![GoDoc](https://godoc.org/github.com/banovo/go-orocrm-wsse?status.svg)](https://godoc.org/github.com/banovo/go-orocrm-wsse) [![Build Status](https://travis-ci.org/banovo/go-orocrm-wsse.svg?branch=master)](https://travis-ci.org/banovo/go-orocrm-wsse) [![Go Report Card](https://goreportcard.com/badge/github.com/banovo/go-orocrm-wsse)](https://goreportcard.com/report/github.com/banovo/go-orocrm-wsse) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md) [![Code Climate](https://codeclimate.com/github/banovo/go-orocrm-wsse/badges/gpa.svg)](https://codeclimate.com/github/banovo/go-orocrm-wsse)
Go WSSE transport (http.RoundTripper)

## Difference to the base package
This is an update of motemen's package. We avoid binary data and changed the date to ISO 8601.
We do this because php and [orocrm](https://www.orocrm.com/documentation/index/current/cookbook/how-to-use-wsse-authentication) and the corresponding package and [EscapeWSSEAuthenticationBundle](https://github.com/djoos/EscapeWSSEAuthenticationBundle) seem to have problems with the orginal data.

## Thanks
Special thanks to [@motemen](https://github.com/motemen) and his [wsse package](https://github.com/motemen/go-wsse) we (needed) to fork from.
