#!/bin/sh
curl http://localhost:8000/createAccount -H "Name:account1" -w %{http_code}
