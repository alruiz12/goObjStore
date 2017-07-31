#!/bin/sh
curl http://localhost:8000/createAccount -H "Nme:account1" -w %{http_code}