#!/bin/sh
curl http://localhost:8000/putObj -T "../bigFile" -w %{http_code}