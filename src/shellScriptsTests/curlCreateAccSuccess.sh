#!/bin/sh
curl http://localhost:8000/putAcc -H "Name:account1" -w %{http_code}
