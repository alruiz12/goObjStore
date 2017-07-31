#!/bin/sh
curl http://localhost:8000/putAcc -H "Nme:account1" -w %{http_code}